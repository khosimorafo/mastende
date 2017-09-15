package tenants

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/now"
	"github.com/khosimorafo/mastende/db"
	"github.com/khosimorafo/mastende/invoices"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/payments"
	"github.com/khosimorafo/mastende/periods"
	"github.com/khosimorafo/mastende/utils"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var a db.App

/**
*
* Create new tenants record.
 */
func New(app *db.App) *Tenant {

	a = *app
	a.SetCollection("tenants")
	t := &Tenant{}
	return t
}

/**
*
* When provided an id, This function returns a tenants with data already read from the database
*
 */

func NewInstanceWithId(app *db.App, id string) (*Tenant, error) {

	a = *app
	tenant := New(&a)
	tenant.ID = id

	if err := tenant.Read(); err != nil {

		return tenant, errors.New("Error while attempting to read tenants.")

	}
	return tenant, nil
}

func (tenant *Tenant) Persist() error {

	tenant.ID = utils.RandStringRunes(25)

	if err := persist(tenant); err != nil {

		return err
	}
	return nil
}

func (tenant *Tenant) Read() error {

	if err := get(tenant); err != nil {

		return errors.New("Error retrieving tenants.")
	}

	payments.ListByTenant(&a, &tenant.Payments, tenant.ID)
	invoices.ListByTenant(&a, &tenant.Invoices, tenant.ID)

	tenant.populateTotals()

	return nil
}

func (tenant *Tenant) Delete() error {

	if err := remove(tenant); err != nil {

		return errors.New("Error retrieving tenants.")
	}
	return nil
}

func (tenant *Tenant) Update() error {

	if err := update(tenant); err != nil {

		return errors.New("Error updating tenants.")
	}
	return nil
}

func (tenant *Tenant) Validate() error {

	// 1. Check if move_in_date is valid.
	_, _, err := utils.DateFormatter(tenant.MoveInDate)

	if err != nil {

		return errors.New("Invalid move_in_date.")
	}

	return nil
}

func (tenant *Tenant) SetStatusActive() error {

	// Setup update
	colQuerier := bson.M{"id": tenant.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Active", "updatedtime": time.Now().String(),
		},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func (tenant *Tenant) SetStatusInActive() error {

	// Setup update
	colQuerier := bson.M{"id": tenant.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "InActive", "updatedtime": time.Now().String(),
		},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func (tenant *Tenant) MonthlyTenantInvoice(input map[string]interface{}) (*invoices.Invoice, error) {

	date := input["date"].(string)
	dayDueBy := input["daydueby"].(int)
	lastDiscountDay := input["lastdiscountday"].(int)

	_, t, err := utils.DateFormatter(date)

	if err != nil {

		log.Fatal("Date parsing error : ", err)
		return nil, err
	}

	dueDate := now.New(t).AddDate(0, 0, dayDueBy).Format("2006-01-02")
	//Set up previous month
	previous := now.New(t).AddDate(0, -1, 0)
	startOfPrevious := now.New(previous).BeginningOfMonth()
	lastDateforDiscount := startOfPrevious.AddDate(0, 0, lastDiscountDay).Format("2006-01-02")

	//log.Println("tenant.go ", lastDateforDiscount)

	period, err := periods.NewInstanceWithDate(&a, date)

	if err != nil {

		return nil, err
	}

	lineItems := tenantItem()

	var (
		total float64 = 0
		lines []invoices.LineItem
	)

	lines = *lineItems

	for _, li := range lines {

		total += li.Rate
	}

	number := fmt.Sprintf("%v-%s-%s", period.Name, tenant.Site, tenant.Room)

	i := invoices.Invoice{

		ID:                  utils.RandStringRunes(25),
		TenantID:            tenant.ID,
		TenantName:          tenant.Name,
		Number:              number,
		Date:                date,
		DueDate:             dueDate,
		LastDateForDiscount: lastDateforDiscount,
		Reference:           "",
		Total:               total,
		Balance:             total,
		Discount:            0,
		LineItems:           lines,
		PeriodIndex:         period.Index,
		PeriodName:          period.Name,
		Status:              "Draft",
	}

	return &i, nil
}

func persist(tenant *Tenant) error {

	// Index
	index := mgo.Index{
		Key:        []string{"zaid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := a.Collection.EnsureIndex(index)

	if err != nil {
		panic(err)
	}

	tenant.UpdateTime = time.Now().String()
	tenant.CreatedTime = time.Now().String()

	if err := a.Collection.Insert(tenant); err != nil {
		log.Println(err)
		return errors.New("Duplicate ZAID submitted.")
	}

	return nil

}

func get(tenant *Tenant) error {

	if err := a.Collection.Find(bson.M{"id": tenant.ID}).One(&tenant); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	} else {

	}

	return nil
}

func remove(tenant *Tenant) error {

	if err := a.Collection.Remove(bson.M{"id": tenant.ID}); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func update(tenant *Tenant) error {

	tenant.UpdateTime = time.Now().String()

	// Setup update
	colQuerier := bson.M{"id": tenant.ID}
	change := bson.M{"$set": bson.M{
		"name":       tenant.Name,
		"zaid":       tenant.ZAID,
		"mobile":     tenant.Mobile,
		"telephone":  tenant.Telephone,
		"imageurl":   tenant.ImageURL,
		"site":       tenant.Site,
		"room":       tenant.Room,
		"gender":     tenant.Gender,
		"status":     tenant.Status,
		"moveindate": tenant.MoveInDate,

		"updatedtime": time.Now().String(),
	},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func TenantList(app *db.App, tens *[]Tenant, input map[string]interface{}) (error) {

	a = *app
	a.SetCollection("tenants")

	list := []Tenant{}

	if err := a.Collection.Find(bson.M{}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*tens = list; return nil
}

func (tenant *Tenant) populateTotals()  {

	invs := tenant.Invoices

	var outstanding float64 = 0.0
	var overdue		float64 = 0.0


	for _, inv := range invs {

		 outstanding += inv.Balance

		if inv.Status == "overdue" {

			overdue += inv.Balance
		}
	}

	tenant.Outstanding = outstanding
	tenant.Overdue =overdue
}

func tenantItem() *[]invoices.LineItem {

	// 1. Get rental item
	item := items.Item{

		ID:          utils.RandStringRunes(25),
		Name:        "",
		Description: "",
		Rate:        330,
		Status:      "Active",
	}

	line := invoices.LineItem{

		InvoiceID:   item.ID,
		Name:        item.Name,
		Description: item.Description,
		Rate:        item.Rate,
		Quantity:    1,
	}

	return &[]invoices.LineItem{line}
}

type Tenant struct {
	ID               string             `json:"id",omitempty"`
	Name             string             `json:"name"`
	ZAID             string             `json:"zaid"`
	Telephone        string             `json:"telephone"`
	Mobile           string             `json:"mobile"`
	Site             string             `json:"site"`
	Room             string             `json:"room"`
	Gender           string             `json:"gender"`
	MoveInDate       string             `json:"moveindate"`
	MoveOutDate      string             `json:"moveoutdate"`
	LastManualPeriod string             `json:"lastmanualperiod"`
	Outstanding      float64            `json:"outstanding"`
	Overdue      	 float64            `json:"overdue"`
	Credits          float64            `json:"credit"`
	Status           string             `json:"status"`
	ImageURL         string             `json:"imageurl,omitempty"`
	Invoices         []invoices.Invoice `json:"invoices,"`
	Payments         []payments.Payment `json:"payments,"`

	CreatedTime string `json:"createdtime"`
	UpdateTime  string `json:"updatedtime"`
}
