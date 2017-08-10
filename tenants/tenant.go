package tenants

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"time"
	"log"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/db"
	"github.com/khosimorafo/mastende/invoices"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/periods"
	"fmt"
	"github.com/khosimorafo/mastende/payments"
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

func (tenant *Tenant) Persist() (error) {

	tenant.ID = utils.RandStringRunes(25)

	if err := persist(tenant); err != nil{

		return err
	}
	return nil
}

func (tenant *Tenant) Read() (error) {

	if err := get(tenant); err != nil{

		return errors.New("Error retrieving tenants.")
	}

	payments.ListByTenant(&a, &tenant.Payments, tenant.ID)
	invoices.ListByTenant(&a, &tenant.Invoices, tenant.ID)

	return nil
}

func (tenant *Tenant) Delete() (error) {

	if err := remove(tenant); err != nil{

		return errors.New("Error retrieving tenants.")
	}
	return nil
}

func (tenant *Tenant) Update() (error) {

	if err := update(tenant); err != nil{

		return errors.New("Error updating tenants.")
	}
	return nil
}

func (tenant *Tenant) Validate() (error) {

	// 1. Check if move_in_date is valid.
	_, _, err := utils.DateFormatter(tenant.MoveInDate)

	if err != nil {

		return errors.New("Invalid move_in_date.")
	}

	return nil
}

func (tenant *Tenant) SetStatusActive() (error) {

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

func (tenant *Tenant) SetStatusInActive() (error) {

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

func (tenant *Tenant) MonthlyInvoice () (*invoices.Invoice, error) {

	period, err := periods.NewInstanceWithDate(&a, tenant.MoveInDate)

	if err != nil { return nil, err }


	lineItems := tenantItem()

	var total float64 = 0

	var lines []invoices.LineItem

	lines = *lineItems

	for _, li := range lines {

		total += li.Rate
	}

	number := fmt.Sprintf("%v-%s-%s", period.Name, tenant.Site, tenant.Room)

	i := invoices.Invoice{

		ID:          utils.RandStringRunes(25),
		TenantID:    tenant.ID,
		TenantName:  tenant.Name,
		Number:      number,
		Date:	     tenant.MoveInDate,
		Reference:   "",
		Total:       total,
		Balance:     total,
		Discount:    0,
		LineItems:   lines,
		PeriodIndex: period.Index,
		PeriodName:  period.Name,
		Status:      "Draft",
	}

	return &i, nil
}

func persist(tenant *Tenant) (error) {

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
	} else{


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

func update(tenant *Tenant) error{

	tenant.UpdateTime = time.Now().String()

	// Setup update
	colQuerier := bson.M{"id": tenant.ID}
	change := bson.M{"$set": bson.M{
					"name": 	tenant.Name,
					"zaid":		tenant.ZAID,
					"mobile":	tenant.Mobile,
					"telephone":	tenant.Telephone,
					"imageurl":	tenant.ImageURL,
					"site":		tenant.Site,
					"room":		tenant.Room,
					"gender":	tenant.Gender,
					"status":	tenant.Status,
					"moveindate":	tenant.MoveInDate,

					"updatedtime": 	time.Now().String(),
					},
			}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func ListTenants (list *[]Tenant) (int8, error) {

	var listSize int8 = 0

	if err := a.Collection.Find(bson.M{}).All(&list); err != nil {

		return listSize, err
	}

	return listSize, nil
}

func tenantItem() *[]invoices.LineItem  {

	// 1. Get rental item
	item := items.Item{

		ID: utils.RandStringRunes(25),
		Name: "",
		Description:"",
		Rate: 330,
		Status: "Active",
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

	ID          		string  		`json:"id",omitempty"`
	Name        		string  		`json:"name"`
	ZAID        		string  		`json:"zaid"`
	Telephone   		string  		`json:"telephone"`
	Mobile      		string  		`json:"mobile"`
	Site        		string  		`json:"site"`
	Room        		string  		`json:"room"`
	Gender        		string  		`json:"gender"`
	MoveInDate 	 	string  		`json:"moveindate"`
	MoveOutDate 		string  		`json:"moveoutdate"`
	LastManualPeriod 	string 			`json:"lastmanualperiod"`
	Outstanding 		float64 		`json:"outstanding"`
	Credits     		float64 		`json:"credit"`
	Status      		string  		`json:"status"`
	ImageURL	   	string  		`json:"imageurl,omitempty"`
	Invoices 		[]invoices.Invoice 	`json:"invoices,"`
	Payments		[]payments.Payment 	`json:"payments,"`

	CreatedTime        	string  		`json:"createdtime"`
	UpdateTime        	string  		`json:"updatedtime"`
}

