package mastende_test

import (
	"testing"
	"os"
	"github.com/khosimorafo/mastende/mastende"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/periods"
	"github.com/khosimorafo/mastende/invoices"
	"encoding/json"
	"fmt"
)


var a db.App
var tenantid string
var invoiceid string

// TestMain wraps all tests with the need fixtures
func TestMain(m *testing.M) {

	// Init test session/db/collection
	configor.Load(&db.Config, "config.yml")
	a = *db.DB(db.Config.DB.DbUrl, db.Config.DB.DbName)
	a.Database.DropDatabase()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// Close session
	a.Session.Close()

	// Call with result of m.Run()
	os.Exit(retCode)
}

/********************************************************************************/
/*********************************Tenant*****************************************/
/********************************************************************************/

func testMastendeCreateTenant(t *testing.T)  {

	testCreatePeriodRange(t)

	id := createTenant(t)
	getTenant(id, t)
	updateTenant(id, t)
	deleteTenant(id, t)
}

func createTenant(t *testing.T)string {

	input := map[string]interface{}{
		"name":       "Mitchell Mastende Test",
		"zaid":       utils.RandStringRunes(13),
		"moveInDate": "2017-08-01",
	}

	_mastende := mastende.New(&a)

	if err := _mastende.CreateTenant(input, true); err != nil {

		t.Errorf("matende.go: Error creating tenants", err.Error())
	}

	return _mastende.Tenant.ID
}

func getTenant(tenantId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   tenantId,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.GetTenant(input); err != nil {

		t.Errorf("matende.go: Error retrieving tenants", err.Error())
	} else {

		t.Log("Tenant Invoices are ", _mastende.Tenant)

		b, err := json.Marshal(_mastende.Tenant)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}
}

func updateTenant(tenantId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   	tenantId,
		"name":   	"Mitchell Updated Tenant",
		"mobile":   	"0823345785",
	}

	_mastende := mastende.New(&a)

	if err := _mastende.UpdateTenant(input); err != nil {

		t.Errorf("mastende.go: Error updating tenants", err.Error())
	}
}

func deleteTenant(tenantId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   	tenantId,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.DeleteTenant(input); err != nil {

		t.Errorf("matende.go: Error deleting tenants", err.Error())
	}
}

/********************************************************************************/
/********************************Invoice*****************************************/
/********************************************************************************/

func TestMastendeInvoice(t *testing.T)  {

	testCreatePeriodRange(t)

	tenantid = createTenant(t)
	//getTenant(tenantid, t)
	updateTenant(tenantid, t)


	id := createInvoice(t)
	getInvoice(id, t)
	updateInvoice(id, t)
	discountInvoice(id, t)
	latePaymentRequest(id, t)
	makeInvoicePayment(id, t)

	getTenant(tenantid, t)

	//testTenantList(t)
	//testInvoiceListByPeriod(t)
	//testOutstandingInvoices(t)

	// deleteInvoice(id, t)
	// deleteTenant(tenantid, t)
}

func createInvoice(t *testing.T)string {

	_mastende := mastende.New(&a)

	// Input for invoice
	input := map[string]interface{}{
		"tenantid":       	tenantid,
		"date":       		"2017-08-01",
		"daydueby":   		4,
		"lastdiscountday":  24,
	}

	if err := _mastende.CreateMonthlyTenantInvoice(input); err != nil {

		t.Errorf("matende.go: Error creating invoices", err.Error())
	}

	return _mastende.Invoice.ID
}

func getInvoice(invoiceId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   invoiceId,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.GetInvoice(input); err != nil {

		t.Errorf("matende.go: Error retrieving invoices", err.Error())
	}
}

func updateInvoice(invoiceId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   invoiceId,
		"date":       		"2017-08-04",
		"duedate": 		"2017-08-10",
	}

	_mastende := mastende.New(&a)

	if err := _mastende.UpdateInvoice(input); err != nil {

		t.Errorf("matende.go: Error updating invoices", err.Error())
	}
}

func deleteInvoice(invoiceId string, t *testing.T) {

	input := map[string]interface{}{
		"id":   	invoiceId,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.DeleteInvoice(input); err != nil {

		t.Errorf("matende.go: Error deleting invoices", err.Error())
	}
}

func discountInvoice(invoiceId string, t *testing.T){

	input := map[string]interface{}{

		"invoiceid": invoiceId,
		"amount": 100.0,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.DiscountInvoice(input); err != nil {

		t.Errorf("matende.go: Error making invoice payment", err.Error())
	}
}

func latePaymentRequest(invoiceId string, t *testing.T){

	input := map[string]interface{}{

		"invoiceid"		: invoiceId,
		"paybydate" 	: "2017-08-20", 
		"requestdate"	: "2017-08-10",
		"requestby" 	: "Khosi Morafo", 
		"requestmode"	: "Cash",
	}

	_mastende := mastende.New(&a)

	if err := _mastende.PaymentExtensionRequest(input); err != nil {

		t.Errorf("mastende.go: Error making late payment request", err.Error())
	}

}

func makeInvoicePayment(invoiceId string, t *testing.T){

	invoice, err := invoices.NewInstanceWithId(&a, invoiceId)
	if err != nil {

		t.Errorf("Error retrieving invoice : ", err.Error())
		return
	}

	input := map[string]interface{}{

		"invoiceid": invoiceId,
		"tenantid": invoice.TenantID,
		"number":	"",
		"description" : "Pay Description",
		"date" : invoice.DueDate, 
		"mode" : "Cash",
		"amount": invoice.Balance,
	}

	_mastende := mastende.New(&a)

	if err := _mastende.MakeInvoicePayment(input); err != nil {

		t.Errorf("matende.go: Error making invoice payment", err.Error())
	}
}

/********************************************************************************/
/********************************Lists*****************************************/
/********************************************************************************/

func TestTenantList(t *testing.T){

	_mastende := mastende.New(&a)

	input := map[string]interface{}{

		"status":       "Active",
	}

	if err := _mastende.TenantList(input); err != nil {

		t.Errorf("matende.go: Error retrieving tenant list - ", err.Error())
	}

	if len(_mastende.Tenants) < 1 {

		t.Errorf("Expected list size, at least,  (%v). Got %v", 1, len(_mastende.Tenants))
	}
}

func TestInvoiceListByPeriod(t *testing.T){

	_mastende := mastende.New(&a)

	input := map[string]interface{}{

		"periodname":       "August-2017",
	}

	if err := _mastende.InvoiceListByPeriod(input); err != nil {

		t.Errorf("mastende.go: Error retrieving invoice list - ", err.Error())
	}

	if len(_mastende.Invoices) < 1 {

		t.Errorf("Expected list size, at least,  (%v). Got %v", 1, len(_mastende.Invoices))
	}
}

func TestOutstandingInvoices(t *testing.T){

	_mastende := mastende.New(&a)

	input := map[string]interface{}{

		"periodname":       "August-2017",
	}

	if err := _mastende.OustandingInvoiceList(input); err != nil {

		t.Errorf("mastende.go: Error retrieving invoice list - ", err.Error())
	}

	if len(_mastende.Invoices) < 1 {

		t.Errorf("Expected list size, at least,  (%v). Got %v", 1, len(_mastende.Invoices))
	}
}

func TestInvoiceListByTenant(t *testing.T){

	_mastende := mastende.New(&a)

	input := map[string]interface{}{

		"tenantid": tenantid,
	}

	if err := _mastende.InvoiceListByTenant(input); err != nil {

		t.Errorf("mastende.go: Error retrieving invoice list - ", err.Error())
	}

	if len(_mastende.Invoices) < 1 {

		t.Errorf("Expected invoice list size, at least,  (%v). Got %v", 1, len(_mastende.Invoices))
	}
}

func TestPaymentListByTenant(t *testing.T){

	_mastende := mastende.New(&a)

	input := map[string]interface{}{

		"tenantid": tenantid,
	}

	if err := _mastende.PaymentListByTenant(input); err != nil {

		t.Errorf("mastende_test.go: Error retrieving payment list - ", err.Error())
	}

	if len(_mastende.Payments) < 1 {

		t.Errorf("Expected payment list size, at least,  (%v). Got %v", 1, len(_mastende.Invoices))
	}
}

func testCreatePeriodRange(t *testing.T)  {

	if err := periods.CreateFinancialPeriodRange(&a, "2017-07-01", 12); err != nil{

		t.Error("Error creating period data!")
		return
	}
}
