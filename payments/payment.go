package payments

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/khosimorafo/mastende/db"
	"errors"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/khosimorafo/mastende/utils"
	"log"
)


var a db.App

/**
*
* Create new payments record.
*/
func New(db *db.App) *Payment {

	a = *db
	a.SetCollection("payments")
	t := &Payment{}
	return t
}

/**
*
* When provided an id, This function returns a payments with data already read from the database
*
*/
func NewInstanceWithId(app *db.App, id string) (*Payment, error) {

	a = *app
	payment := New(&a)
	payment.ID = id

	if err := payment.Read(); err != nil {

		return payment, errors.New("Error while attempting to read payments.")

	}
	return payment, nil
}

func (payment *Payment) Persist() (error) {

	payment.ID = utils.RandStringRunes(25)

	if err := persist(payment); err != nil{

		return errors.New("Error while attempting to save payments.")
	}
	return nil
}

func (payment *Payment) Read() (error) {

	if err := get(payment); err != nil{

		return errors.New("Error retrieving payments.")
	}
	return nil
}

func (payment *Payment) Delete() (error) {

	if err := delete(payment); err != nil{

		return errors.New("Error retrieving payments.")
	}
	return nil
}

func (payment *Payment) Update() (error) {

	if err := update(payment); err != nil{

		return errors.New("Error updating payments.")
	}
	return nil
}

func (payment *Payment) Validate() (error) {

	// 1. Check if move_in_date is valid.
	_, _, err := utils.DateFormatter(payment.Date)

	if err != nil {

		return errors.New("Invalid payment date.")
	}

	return nil
}

func (payment *Payment) SetStatusActive() (error) {

	// Setup update
	colQuerier := bson.M{"id": payment.ID}
	change := bson.M{
		"$set": bson.M{

			"status": "Active",
			"updatedtime": time.Now().String(),
		},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func (payment *Payment) SetStatusInActive() (error) {

	// Setup update
	colQuerier := bson.M{"id": payment.ID}
	change := bson.M{
		"$set": bson.M{

			"status": "InActive",
			"updatedtime": time.Now().String(),
		},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func persist(payment *Payment) (error) {

	// Index
	index := mgo.Index{
		Key:        []string{"invoiceid", "customerid"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}

	err := a.Collection.EnsureIndex(index)

	if err != nil {
		panic(err)
	}

	if err := a.Collection.Insert(payment); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func get(payment *Payment) error {

	if err := a.Collection.Find(bson.M{"id": payment.ID}).One(&payment); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func delete(payment *Payment) error {

	if err := a.Collection.Remove(bson.M{"id": payment.ID}); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func update(payment *Payment) error{

	// Setup update
	colQuerier := bson.M{"id": payment.ID}
	change := bson.M{"$set": bson.M{

			"number": 		payment.Number,
			"description":		payment.Description,
			"invoiceid":		payment.InvoiceID,
			"tenantid":		payment.TenantID,
			"mode":			payment.Mode,
			"amount":		payment.Amount,
			"status":		payment.Status,
			"date":			payment.Date,
			"updatedtime": 		time.Now().String(),
		},
	}

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func ListAllPayments (app *db.App, payments *[]Payment) (error) {

	a = *app
	a.SetCollection("payments")

	list := []Payment{}

	if err := a.Collection.Find(bson.M{}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*payments = list; return nil
}

func ListPaymentsByInvoice (app *db.App, payments *[]Payment, invoiceid string) (error) {

	a = *app
	a.SetCollection("payments")

	list := []Payment{}

	if err := a.Collection.Find(bson.M{"invoiceid": invoiceid}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*payments = list; return nil
}

func ListByTenant(app *db.App, payments *[]Payment, tenantid string) (error) {

	a = *app
	a.SetCollection("payments")

	list := []Payment{}

	if err := a.Collection.Find(bson.M{"tenantid": tenantid}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*payments = list; return nil
}

type Payment struct {

	ID          	string   `json:"id,omitempty"`
	Description 	string   `json:"description"`
	TenantID    	string   `json:"tenantid"`
	InvoiceID   	string   `json:"invoiceid"`
	Number      	string   `json:"number"`
	Amount      	float64  `json:"amount"`
	Mode        	string   `json:"mode"`
	Date        	string   `json:"date"`
	Status      	string   `json:"status"`

	CreatedTime	string   `json:"createdtime"`
	UpdatedTime 	string   `json:"updatedtime"`
}
