package invoices

import (
	"github.com/khosimorafo/mastende/db"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/payments"
	"github.com/khosimorafo/mastende/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
	"fmt"
	"log"
	"time"
)

var a db.App

/**
*
* Create new invoices record.
 */
func New(db *db.App) *Invoice {

	a = *db
	a.SetCollection("invoices")
	t := &Invoice{}
	return t
}

/**
*
* When provided an id, This function returns a invoice with data already read from the database.
*
 */
func NewInstanceWithId(app *db.App, id string) (*Invoice, error) {

	a = *app
	invoice := New(&a)
	invoice.ID = id

	if err := invoice.Read(); err != nil {

		return invoice, errors.New("Error while attempting to read invoice.")

	}
	return invoice, nil
}

/**
*
* Provides scaffolding for a new tenant.
*
 */
func NewInstanceOfTenantInvoice(app *db.App, input map[string]interface{}) (*Invoice, error) {

	a = *app
	invoice := newTenantInvoice(input)

	if err := invoice.Persist(); err != nil {

		return invoice, errors.New("Error while attempting to create invoice.")

	}
	return invoice, nil
}

func (invoice *Invoice) Persist() error {

	invoice.ID = utils.RandStringRunes(25)

	if err := persist(invoice); err != nil {

		return errors.New("Error while attempting to save invoices.")
	}

	invoice.UpdateTotalBalanceAndStatus()

	return nil
}

func (invoice *Invoice) Read() error {

	if err := get(invoice); err != nil {

		return errors.New("Error retrieving invoices.")
	}
	return nil
}

func (invoice *Invoice) Delete() error {

	if err := delete(invoice); err != nil {

		return errors.New("Error retrieving invoices.")
	}
	return nil
}

func (invoice *Invoice) Update() error {

	if err := update(invoice); err != nil {

		return errors.New("Error updating invoices.")
	}
	return nil
}

func (invoice *Invoice) Validate() error {

	return nil
}

func (invoice *Invoice) ApplyDiscount(input map[string]interface{}) error {

	var (
		_discount_amount float64 = 0.0
	)
	_discount_amount = input["amount"].(float64)

	if invoice.Total < _discount_amount {

		return errors.New("Discount amount cannot be bigger than invoice total amount.")
	}

	/**
	*  This ensures that only non negative invoice balance amount are accepted.
	 */
	//if invoice.Balance < _discount_amount {
	//
	//	return errors.New("Discount amount cannot be bigger than invoice balance amount.")
	//}

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"discount":    _discount_amount,
			"updatedtime": time.Now().String(),
		},
	}

	if err := db.Update(&a, colQuerier, change); err != nil {

		return err
	}

	invoice.UpdateTotalBalanceAndStatus()

	return nil
}

func (invoice *Invoice) AddLatePaymentFine(app *db.App) error {

	if err := invoice.Read(); err != nil {

		return err
	}

	var request map[string]string
	var payByDate string

	// If a fine addition is attempted and a pay extension exists,
	// Raise error.
	if len(invoice.ExtensionRequest) > 0 {

		request = invoice.ExtensionRequest[0]

		payByDate = request["paybydate"]
		_, d, _ := utils.DateFormatter(payByDate)

		// Test payByDate against current time.
		if time.Now().Before(d) {

			err_str := fmt.Sprintf("Payment request in place. Expires on % ", payByDate)
			return errors.New(err_str)
		}
	}

	// Retrieve Item
	var i []items.Item

	err := items.ItemListByName(app, &i, "Late Payment Fine")

	if err != nil {

		return err
	}

	fmt.Printf("\nItem are ", i)

	item := i[0]

	var (
		_quantity float64 = 1
	)

	lineMap := map[string]interface{}{

		"itemid":      utils.RandStringRunes(25),
		"invoiceid":   invoice.ID,
		"name":        item.Name,
		"description": item.Description,
		"quantity":    _quantity,
		"rate":        item.Rate,
	}

	if e := invoice.AddLineItem(lineMap); e != nil {

		return e
	}

	invoice.UpdateTotalBalanceAndStatus()

	return nil
}

func (invoice *Invoice) AddLineItem(input map[string]interface{}) error {

	line := LineItem{

		ItemID:      input["itemid"].(string),
		InvoiceID:   input["invoiceid"].(string),
		Name:        input["name"].(string),
		Description: input["description"].(string),
		Quantity:    input["quantity"].(float64),
		Rate:        input["rate"].(float64),
	}

	var lines []LineItem
	var _total float64 = 0.0

	lines = invoice.LineItems
	lines = append(lines, line)

	for _, li := range lines {

		_total += li.Rate
	}

	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"total":       _total,
			"lineitems":   lines,
			"updatedtime": time.Now().String(),
		},
	}

	if err := db.Update(&a, colQuerier, change); err != nil {

		return err
	}

	invoice.UpdateTotalBalanceAndStatus()

	return nil
}

func (invoice *Invoice) MakePayment(input map[string]interface{}) error {

	var _amount float64 = 0.0

	_amount = input["amount"].(float64)

	if invoice.Balance > _amount {

		errors.New("Payment amount may not be larger than the balance!")
	}

	payment := payments.New(&a)

	payment.InvoiceID = invoice.ID
	payment.TenantID = invoice.TenantID
	payment.Description = input["description"].(string)
	payment.Amount = input["amount"].(float64)
	payment.Date = input["date"].(string)
	payment.Mode = input["mode"].(string)

	payment.CreatedTime = time.Now().String()
	payment.UpdatedTime = time.Now().String()

	if err := payment.Persist(); err != nil {

		errors.New("Error creating payment!")
	}

	invoice.AddPaymentLine(input)

	return nil
}

func (invoice *Invoice) AddPaymentLine(input map[string]interface{}) error {

	payment := payments.Payment{

		InvoiceID: input["invoiceid"].(string),
		TenantID:  input["tenantid"].(string),
		//Description: input["description"].(string),
		Date:   input["date"].(string),
		Amount: input["amount"].(float64),
		Mode:   input["mode"].(string),
	}

	var payments []payments.Payment

	payments = invoice.Payments
	payments = append(payments, payment)

	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"payments":    payments,
			"updatedtime": time.Now().String(),
		},
	}

	if err := db.Update(&a, colQuerier, change); err != nil {

		return err
	}

	invoice.UpdateTotalBalanceAndStatus()

	return nil
}

func (invoice *Invoice) PaymentExtensionRequest(input map[string]interface{}) error {

	if err := invoice.Read(); err != nil {

		return err
	}

	var (

		requestDate     string = ""
		payByDate 		string = ""
		requestBy       string = ""
		requestMode     string = "Cash"
	)



	if err := utils.DateResolverDefaultOnToday(&requestDate, input["requestdate"]); err != nil {

		return errors.New("Invalid Request Date. ")
	}

	if err := utils.DateResolverDefaultOnToday(&payByDate, input["paybydate"]); err != nil {

		return errors.New("Invalid Pay Date. ")
	}
	
	utils.StringResolver(&requestBy, input["requestby"])
 	utils.StringResolver(&requestMode, input["requestmode"])
	
	request := map[string]string{

		"invoiceid":   invoice.ID,
		"paybydate":   payByDate,
		"requestby":   requestBy,
		"requestdate": requestDate,
		"requestmode": requestMode,
	}

	var requests []map[string]string

	requests = invoice.ExtensionRequest
	requests = append(requests, request)

	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"extensionrequest": requests,
			"updatedtime":      time.Now().String(),
		},
	}

	if err := db.Update(&a, colQuerier, change); err != nil {

		return err
	}

	return nil
}

func (invoice *Invoice) UpdateTotalBalanceAndStatus() error {

	if err := invoice.Read(); err != nil {

		return err
	}

	var (
		balance       float64 = 0.0
		invoice_total float64 = 0.0
		paid_amount   float64 = 0.0
		status        string
	)

	// Source amount already paid
	paid_amount = invoice.GetPaidAmount()
	invoice_total = invoice.GetTotalAmount()

	// Update invoice balance
	balance = invoice_total - (invoice.Discount + paid_amount)

	if invoice_total != 0.0 {

		if balance == 0.0 {

			status = "Paid"
		} else if balance == invoice.Total {

			status = "Unpaid"
		} else if balance < invoice.Total {

			status = "Partial"
		}

		if status != "Paid" {

			_, d, _ := utils.DateFormatter(invoice.DueDate)

			if time.Now().After(d) {

				status = "Overdue"
			}
		}
	}

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"total":       invoice_total,
			"balance":     balance,
			"status":      status,
			"updatedtime": time.Now().String(),
		},
	}

	if err := db.Update(&a, colQuerier, change); err != nil {

		return errors.New("Failed to update total and balance!")
	}

	return nil
}

func (invoice *Invoice) GetPaidAmount() float64 {

	var _paid_amount float64 = 0.0

	for _, payment := range invoice.Payments {

		_paid_amount += payment.Amount
	}
	return _paid_amount
}

func (invoice *Invoice) GetTotalAmount() float64 {

	var _amount float64 = 0.0

	for _, item := range invoice.LineItems {

		_amount += item.Rate
	}
	return _amount
}

func (invoice *Invoice) SetStatusAsDraft() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Draft", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) SetStatusAsSent() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Sent", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) SetStatusAsPartial() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Partial", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) SetStatusAsOverdue() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Overdue", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) SetStatusAsUnpaid() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Partial", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) SetStatusAsPaid() error {

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{
		"$set": bson.M{
			"status": "Paid", "updatedtime": time.Now().String(),
		},
	}

	return db.Update(&a, colQuerier, change)
}

func (invoice *Invoice) GetAsMap(m map[string]interface{}) error {

	if err := get(invoice); err != nil {

		return errors.New("Error creating map.")
	}
	return nil
}

func persist(invoice *Invoice) error {

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

	invoice.UpdatedTime = time.Now().String()
	invoice.CreatedTime = time.Now().String()

	if err := a.Collection.Insert(invoice); err != nil {
		log.Println(err)
		return errors.New("Duplicate ZAID submitted.")
	}
	return nil
}

func get(invoice *Invoice) error {

	if err := a.Collection.Find(bson.M{"id": invoice.ID}).One(&invoice); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func delete(invoice *Invoice) error {

	if err := a.Collection.Remove(bson.M{"id": invoice.ID}); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func update(invoice *Invoice) error {

	invoice.UpdatedTime = time.Now().String()

	// Setup update
	colQuerier := bson.M{"id": invoice.ID}
	change := bson.M{"$set": bson.M{
		"tenantid":    invoice.TenantID,
		"tenantname":  invoice.TenantName,
		"number":      invoice.Number,
		"reference":   invoice.Reference,
		"total":       invoice.Total,
		"balance":     invoice.Balance,
		"discount":    invoice.Discount,
		"date":        invoice.Date,
		"duedate":     invoice.DueDate,
		"periodindex": invoice.PeriodIndex,
		"periodname":  invoice.PeriodName,
		"status":      invoice.Status,
		"updatedtime": invoice.UpdatedTime,
	},
	}

	return db.Update(&a, colQuerier, change)
}

func newTenantInvoice(input map[string]interface{}) *Invoice {

	lineItems := tenantItem()

	var total float64 = 0

	var lines []LineItem

	lines = *lineItems

	for _, li := range lines {

		total += li.Rate
	}

	return &Invoice{

		ID:          utils.RandStringRunes(25),
		TenantID:    input["tenantid"].(string),
		TenantName:  input["tenantname"].(string),
		Number:      "",
		Reference:   "",
		Total:       total,
		Balance:     total,
		Discount:    0,
		LineItems:   lines,
		PeriodIndex: 0,
		PeriodName:  "",
		Status:      "Draft",
	}
}

func tenantItem() *[]LineItem {

	// 1. Get rental item
	item := items.Item{
		ID:          utils.RandStringRunes(25),
		Name:        "",
		Description: "",
		Rate:        330,
		Status:      "Active",
	}

	line := LineItem{

		InvoiceID:   item.ID,
		Name:        item.Name,
		Description: item.Description,
		Rate:        item.Rate,
		Quantity:    1,
	}

	return &[]LineItem{line}
}

func ListByTenant(app *db.App, invoices *[]Invoice, tenantid string) error {

	a = *app
	a.SetCollection("invoices")

	list := []Invoice{}

	if err := a.Collection.Find(bson.M{"tenantid": tenantid}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*invoices = list; return nil
}

func ListByPeriod(app *db.App, invoices *[]Invoice, period_name string) error {

	a = *app
	a.SetCollection("invoices")

	list := []Invoice{}

	if err := a.Collection.Find(bson.M{"periodname": period_name}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*invoices = list; return nil
}

func ListByStatus(app *db.App, invoices *[]Invoice, status string) error {

	a = *app
	a.SetCollection("invoices")

	list := []Invoice{}

	if err := a.Collection.Find(bson.M{"status": status}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*invoices = list; return nil
}

func ListOutstanding(app *db.App, invoices *[]Invoice, period_name string) error {

	a = *app
	a.SetCollection("invoices")

	list := []Invoice{}

	toQuery := []string{"Draft", "Sent", "Partial", "Overdue"}

	if err := a.Collection.Find(bson.M{"status": bson.M{"$in": toQuery}}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the payments and return
	*invoices = list; return nil
}

type Invoice struct {
	ID                  string             `json:"id,omitempty"`
	TenantID            string             `json:"tenantid"`
	TenantName          string             `json:"tenantname"`
	Number              string             `json:"number"`
	Reference           string             `json:"reference"`
	Total               float64            `json:"total"`
	Balance             float64            `json:"balance"`
	Discount            float64            `json:"discount"`
	Date                string             `json:"date"`
	DueDate             string             `json:"duedate"`
	LastDateForDiscount string             `json:"lastdatefordiscount"`
	LineItems           []LineItem         `json:"lineitems"`
	Payments            []payments.Payment `json:"payments"`
	PeriodIndex         int                `json:"periodindex"`
	PeriodName          string             `json:"periodname"`
	Status              string             `json:"status"`

	ExtensionRequest []map[string]string `json:"extension"`

	CreatedTime string `json:"createdtime"`
	UpdatedTime string `json:"updatedtime"`
}
