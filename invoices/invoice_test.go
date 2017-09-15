package invoices_test

import (
	"testing"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
	"os"
	"fmt"
	"github.com/khosimorafo/mastende/invoices"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/tenants"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/periods"
)

var a db.App
var period_name string = ""


// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {

	// Init test session/db/collection
	configor.Load(&db.Config, "config.yml")
	a = *db.DB(db.Config.DB.DbUrl, db.Config.DB.DbName)
	a.Database.DropDatabase()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	a.Session.Close()

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestInvoiceCreate(t *testing.T) {

	// Pass collection to package

	testCreatePeriodRange(t)

	invoiceID := testCreate(t)
	testGet(invoiceID, t)
	testUpdate(invoiceID, t)

	//testGet(invoiceID, t)
	testAddLineItem(invoiceID, t)
	testPaymentExtensionRequest(invoiceID, t)
	testAddLatePaymentFine(invoiceID, t)
	testApplyDiscount(invoiceID, t)

	//testMakePayment(invoiceID, t)
	//testDelete(invoiceID, t)
}

func TestListByPeriod(t *testing.T){

	list := []invoices.Invoice{}

	if err := invoices.ListByPeriod(&a, &list, period_name); err != nil {

		t.Errorf("Error creating tenants", err.Error())
	}

	if len(list) < 1 {

		t.Errorf("Expected list size of (%v). Got %v", 1, len(list))
	}
}


func TestListOutstanding(t *testing.T){

	list := []invoices.Invoice{}

	if err := invoices.ListOutstanding(&a, &list, period_name); err != nil {

		t.Errorf("Error creating tenants", err.Error())
	}

	if len(list) < 1 {

		t.Errorf("Expected list size of (%v). Got %v", 1, len(list))
	}
}

func testCreate(t *testing.T) string {

	tenant := tenants.New(&a)
	tenant.Name = "Test Tenant Name Name"
	tenant.ZAID = utils.RandNumberRunes(13)

	if err := tenant.Persist(); err != nil {
		t.Errorf("Error creating tenants", err.Error())
	}

	invoice := invoices.New(&a)
	invoice.TenantID = tenant.ID
	invoice.TenantName = tenant.Name
	invoice.Date = "2017-08-01"
	invoice.DueDate = "2017-08-05"

	period, err := periods.NewInstanceWithDate(&a, invoice.Date)

	if err != nil {

		t.Error("Error in assigning period data : ", err.Error())
		return ""
	}

	invoice.PeriodIndex = period.Index
	invoice.PeriodName = period.Name

	if err := invoice.Persist(); err != nil {

		t.Errorf("Error creating invoices", err.Error())
	}

	return invoice.ID
}

func testGet(invoiceId string, t *testing.T) {

	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {
		t.Errorf("Error retrieving invoices : ", err.Error())
	}

	period_name = invoice.PeriodName
	fmt.Printf("Retrieved invoices is : ", invoice)
}

func testUpdate(invoiceId string, t *testing.T) {

	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {
		t.Errorf("Error retrieving invoices : ", err.Error())
	} 

	update_name := "Updated Name "
	invoice.TenantName = update_name
	invoice.Update()

	invoice.Read()

	if invoice.TenantName != update_name {
		t.Errorf("Expected the name to change to (%v). Got %v", update_name, invoice.TenantName)
	}
}

func testDelete(invoiceId string, t *testing.T) {

	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {
		t.Errorf("Error getting invoices for delete : ", err.Error())
	}

	tenantId := invoice.TenantID

	if err = invoice.Delete(); err != nil {
		t.Errorf("Error deleting invoices", err.Error())
		return
	}

	// Remove Tenant
	tenant, err := tenants.NewInstanceWithId(&a, tenantId)
	if err != nil {
		t.Errorf("Error getting tenants for delete : ", err.Error())
	}

	if err = tenant.Delete(); err != nil {
		t.Errorf("Error deleting tenants", err.Error())
		return
	}
}

func testAddLineItem(invoiceId string, t *testing.T) {


	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {
		t.Errorf("Error retrieving invoices : ", err.Error())
	}

	var (
		_quantity	float64 = 1
		_rate 		float64 = 330.0
	)

	line := map[string]interface{}{

		"invoiceid": invoiceId,
		"itemid": utils.RandStringRunes(25),
		"name": "Test Line Item",
		"description" : "Test Line Item Description",
		"quantity": _quantity,
		"rate": _rate,
	}

	if err := invoice.AddLineItem(line); err != nil{

		t.Errorf("Error adding invoice lineitem", err.Error())
		return
	}
}

func testPaymentExtensionRequest(invoiceId string, t *testing.T)  {
	
	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {

		t.Errorf("Error retrieving invoice : ", err.Error())
		return
	}

	input := map[string]interface{}{

		"invoiceid"		: invoice.ID,
		"paybydate" 	: "2017-08-20", 
		"requestdate"	: "2017-08-10",
		"requestby" 	: "Khosi Morafo", 
		"requestmode"	: "Cash",
	}

	if err := invoice.PaymentExtensionRequest(input); err != nil{

		t.Errorf("Error making payment extension request", err.Error())
		return
	}
}

func testAddLatePaymentFine(invoiceId string, t *testing.T) {

	//Add late payment item
	item := items.New(&a)
	item.Name = "Late Payment Fine"
	item.Description = "Payment Fine Description"

	item.Rate = 50

	if err := item.Persist(); err != nil {
		t.Errorf("Error creating items", err.Error())
	}



	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {

		t.Errorf("Error retrieving invoice : ", err.Error())
		return
	}

	if err := invoice.AddLatePaymentFine(&a); err != nil{

		t.Errorf("Error adding late payment fine", err.Error())
		return
	}
}

func testApplyDiscount (invoiceId string, t *testing.T){

	input := map[string]interface{}{

		"invoiceid": invoiceId,
		"amount": 100.0,
	}


	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {

		t.Errorf("Error retrieving invoice : ", err.Error())
		return
	}

	if err := invoice.ApplyDiscount(input); err != nil{

		t.Errorf("Error applying discount", err.Error())
		return
	}
}

func testMakePayment (invoiceId string, t *testing.T)  {

	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {

		t.Errorf("Error retrieving invoice : ", err.Error())
		return
	}

	input := map[string]interface{}{

		"invoiceid": invoice.ID,
		"tenantid" : invoice.TenantID,
		"description" : "Pay Description",
		"date" : "2017-08-01", 
		"mode" : "Cash",
		"amount": invoice.Balance,
	}

	if err := invoice.MakePayment(input); err != nil{

		t.Errorf("Error making payment", err.Error())
		return
	}
}

func testCreatePeriodRange(t *testing.T)  {

	if err := periods.CreateFinancialPeriodRange(&a, "2017-07-01", 12); err != nil{

		t.Error("Error creating period data!")
		return
	}
}