package periods

import (
	"time"
	"github.com/dariubs/percent"
	"github.com/jinzhu/now"
	"github.com/aodin/date"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/khosimorafo/mastende/db"
	"errors"
	"github.com/khosimorafo/mastende/utils"
)

var a db.App

/**
*
* Create new items record.
*/
func New(app *db.App) *Period {

	a = *app
	a.SetCollection("periods")
	t := &Period{}
	return t
}

/***
*
* When provided an id, This function returns a period with data already read from the database
*
*/
func NewInstanceWithIndex(app *db.App, index int) (*Period, error) {

	a = *app
	period := New(&a)

	if err := GetByIndex(period, index); err != nil {

		return period, errors.New("Error while attempting to retrieve period.")

	}
	return period, nil
}

/***
*
* When provided an id, This function returns a period with data already read from the database
*
*/
func NewInstanceWithName(app *db.App, name string) (*Period, error) {

	a = *app
	period := New(&a)

	if err := GetByName(period, name); err != nil {

		return period, errors.New("Error while attempting to retrieve period.")

	}
	return period, nil
}

/***
*
* When provided an id, This function returns a period with data already read from the database
*
*/
func NewInstanceWithDate(app *db.App, date_string string) (*Period, error) {

	a = *app
	period := New(&a)

	if err := GetByDate(period, date_string); err != nil {

		return period, errors.New("Error while attempting to retrieve period.")

	}
	return period, nil
}


type PeriodInterface interface {

	CreateFinancialPeriodRange (start_date string, no_of_months int) (error)
	ReadFinancialPeriodRange (status string) ([]Period, error)
}

func CreateFinancialPeriodRange (app *db.App, start_date string, no_of_months int) (error) {

	app.SetCollection("periods")

	t, err := now.Parse(start_date)

	if err != nil {

		log.Fatal("Date parsing error : ", err)
		return err
	}

	for i := 0; i < no_of_months; i++ {

		current := now.New(t).AddDate(0, i, 0)

		//t.Format(time.RFC3339)
		//current := t.Format("2006-01-02")

		start := now.New(current).BeginningOfMonth().Format("2006-01-02")
		end := now.New(current).EndOfMonth().Format("2006-01-02")

		month := now.New(current).Month()
		year := now.New(current).Year()

		name := fmt.Sprintf("%s-%d", month, year)

		period := Period{i, name, "open", start,end, year, int(month)}

		app.Collection.Insert(period)

	}

	return nil
}

func GetRange (app *db.App, status string) ([]Period, error) {

	app.SetCollection("periods")

	ps := []Period{}
	err := app.Collection.Find(bson.M{}).All(&ps)

	if err != nil {

		return nil, err
	}

	return ps, nil
}

func ReadFinancialPeriodRange (status string) ([]Period, error) {

	ps := []Period{}
	err := a.Collection.Find(bson.M{}).All(&ps)

	if err != nil {

		return nil, err
	}

	return ps, nil
}

func (p *P) GetProRataDays() (float64, error)  {

	days, all, err := p.GetDaysLeft()

	if err != nil {

		return -1, err
	}

	perc := percent.PercentOf(days, all)

	return perc/100, nil
}

func (p *P) GetDaysLeft() (int, int,error)  {

	period, err := p.GetPeriod()

	if err != nil {

		return -1, -1, err
	}

	end, err1 := now.Parse(period.End)
	if err1 != nil {

		return -1, -1,  err
	}

	var no_of_days date.Range
	no_of_days.Start = date.New(p.Date.Date())
	no_of_days.End = date.New(end.Date())

	start, err2 := now.Parse(period.Start)
	if err2 != nil {

		return -1, -1,  err2
	}
	var days_in_month date.Range
	days_in_month.Start = date.New(start.Date())
	days_in_month.End = date.New(end.Date())

	return no_of_days.Days(), days_in_month.Days(), nil

}

func (p *P) GetPeriod () (Period, error) {

	actual_date := date.New(p.Date.Date())

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	for _, period := range ps {

		p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if actual_date.Within(p_range){

			return period, nil
		}
	}

	return Period{}, nil
}

func GetByName (period *Period, name string) (error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return err
	}

	for _, p := range ps {

		if p.Name == name{

			*period = p
		}
	}
	return nil
}

func GetByIndex (period *Period, index int) (error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return err
	}

	for _, p := range ps {

		if p.Index == index{

			*period = p
		}
	}
	return nil
}

func GetByDate (period *Period, str_date string) (error) {

	_, d, err := utils.DateFormatter(str_date)

	if err != nil {

		return err
	}

	err = nil

	p := P{ Date: d}
	ps, err := p.GetPeriod()

	if err != nil {

		return err
	}

	*period = ps

	return nil
}

func GetLatestPeriod() (*Period, error){

	//current_date, _, _ := DateFormatter(time.Now().UTC().String())

	actual_date := date.New(time.Now().UTC().Date())

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return &Period{}, err
	}

	for _, period := range ps {

		p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if actual_date.Within(p_range){

			return &period, nil
		}
	}

	return &Period{}, nil
}

func GetSequentialPeriodRange(start string, end string) ([]Period, error){

	ps, err := ReadFinancialPeriodRange("open")

	p := New(&a)
	p_latest := New(&a)

	GetByDate(p, start)
	GetByDate(p, end)


	if err != nil {

		return []Period{}, err
	}

	periods := make([]Period, 0)
	for _, period := range ps {

		if (period.Index >= p.Index) && (period.Index <= p_latest.Index) {

			periods = append(periods, period)
		}
	}
	return periods, nil
}

func GetSequentialPeriodRangeFromToCurrent(start string) ([]Period, error){

	ps, err := ReadFinancialPeriodRange("open")

	p := New(&a)
	p_latest := New(&a)

	GetByDate(p, start)
	p_latest, _ = GetLatestPeriod()


	if err != nil {

		return []Period{}, err
	}

	periods := make([]Period, 0)
	for _, period := range ps {

		if (period.Index >= p.Index) && (period.Index <= p_latest.Index) {

			periods = append(periods, period)
		}
	}
	return periods, nil
}

func GetSequentialPeriodRangeAfterToCurrent(start string) ([]Period, error){

	ps, err := ReadFinancialPeriodRange("open")
	p := New(&a)

	GetByDate(p, start)
	p_latest,err := GetLatestPeriod()


	if err != nil {

		return []Period{}, err
	}

	periods := make([]Period, 0)
	for _, period := range ps {

		if (period.Index > p.Index) && (period.Index <= p_latest.Index) {

			periods = append(periods, period)
		}
	}
	return periods, nil
}

func (period Period) GetPeriodDiscountDate() (time.Time, bool)  {

	_, p_start_t, _ := utils.DateFormatter(period.Start)

	//Create a stub(holder) date that navigates to the previous month.
	d := time.Duration(-int(p_start_t.Day())-5) * 24 * time.Hour
	stub_date := p_start_t.Add(d)

	//Go to the beginning of the previous month and add 25 day. The result is the cut-off date/time.
	d_date := now.New(stub_date).BeginningOfMonth().AddDate(0,0,25)

	//Use cut-off date/time to create a before range
	beforeCutoff := date.Range{End: date.New(d_date.Year(), d_date.Month(), d_date.Day())}

	//Test against today. Assumes that today's date/time is the actual test date.
	today := date.FromTime(time.Now())
	if (today.Within(beforeCutoff)){

		return d_date, true
	}

	return d_date, false
}

func GetNextPeriodByName (name string) (Period, error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	var isnext bool
	isnext = false
	for _, period := range ps {

		if isnext {

			return period, nil
		}
		//p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if period.Name == name{

			isnext = true
		}
	}

	return Period{}, nil
}

type P struct {

	Date time.Time
}

type Period struct {

	Index int 	`json:"index,omitempty"`
	Name string 	`json:"name,omitempty"`
	Status string 	`json:"status,omitempty"`

	Start string 	`json:"start_date,omitempty"`
	End string 	`json:"end_date,omitempty"`
	Year int	`json:"year,omitempty"`
	Month int	`json:"month,omitempty"`
}
