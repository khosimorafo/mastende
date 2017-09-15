package payments_test

import (
	"os"
	"testing"
	"fmt"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
	"github.com/khosimorafo/mastende/payments"
	"time"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/periods"
	"github.com/khosimorafo/mastende/tenants"
	"github.com/khosimorafo/mastende/invoices"
)

var a db.App

var tenantid string
var invoiceid string

// TestMain wraps all tests with the needed initialized mock DB and fixtures
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

func testCreateTenantAndInvoice(t *testing.T) {

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
	}

	invoice.PeriodIndex = period.Index
	invoice.PeriodName = period.Name

	if err := invoice.Persist(); err != nil {

		t.Errorf("Error creating invoices", err.Error())
	}

	tenantid = tenant.ID
	invoiceid = invoice.ID
}

func TestPaymentCreate(t *testing.T) {

	// Pass collection to package
	testCreateTenantAndInvoice(t)

	paymentID := testCreate(t)
	testGet(paymentID, t)
	testUpdate(paymentID, t)
	testGet(paymentID, t)
	//testDelete(paymentID, t)

}

func testCreate(t *testing.T) string {

	payment := payments.New(&a)

	_date, _, _ := utils.DateFormatter(time.Now().String())
	payment.TenantID = tenantid
	payment.InvoiceID = invoiceid
	payment.Date =  _date
	payment.Mode = "Cash"
	payment.Amount = 333.0

	if err := payment.Persist(); err != nil {
		t.Errorf("Error creating payments", err.Error())
	}

	return payment.ID
}

func testGet(paymentId string, t *testing.T) {

	payment, err := payments.NewInstanceWithId(&a, paymentId)
	if err != nil {
		t.Errorf("Error retrieving payments : ", err.Error())
	}

	fmt.Printf("Retrieved payments is : ", payment)
}

func testUpdate(paymentId string, t *testing.T) {

	payment, err := payments.NewInstanceWithId(&a, paymentId)
	if err != nil {
		t.Errorf("Error retrieving payments : ", err.Error())
	}

	_date, _, _ := utils.DateFormatter(time.Now().String())

	payment.Date =  _date
	payment.Mode = "Cash"
	payment.Amount = 333.0
	payment.Update()
}

func testDelete(paymentId string, t *testing.T) {

	payment, err := payments.NewInstanceWithId(&a, paymentId)
	if err != nil {
		t.Errorf("Error getting payments for delete : ", err.Error())
	}

	if err = payment.Delete(); err != nil {
		t.Errorf("Error deleting payments", err.Error())
		return
	}
}


func TestListByTenant(t *testing.T) {

	list := []payments.Payment{}

	if err := payments.ListByTenant(&a, &list, tenantid); err != nil {

		t.Errorf("Error creating tenants", err.Error())
	}

	if len(list) < 1 {

		t.Errorf("Expected list size of (%v). Got %v", 1, len(list))
	}
}
