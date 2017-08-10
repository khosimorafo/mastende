package mastende_test

import (
	"testing"
	"os"
	"github.com/khosimorafo/mastende/mastende"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/periods"
)

var a db.App

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

func TestMastendeCreateTenant(t *testing.T)  {

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

		t.Errorf("matende.go: Error updating tenants", err.Error())
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

	id := createInvoice(t)
	getInvoice(id, t)
	updateInvoice(id, t)
	//deleteInvoice(id, t)
}

func createInvoice(t *testing.T)string {

	_mastende := mastende.New(&a)

	// Input for tenant
	ten_input := map[string]interface{}{
		"name":       "Mitchell Mastende Test",
		"zaid":       utils.RandStringRunes(13),
		"moveindate": "2017-08-01",
	}

	if err := _mastende.CreateTenant(ten_input, false); err != nil {

		t.Errorf("matende.go: Error creating tenants", err.Error())
	}

	// Input for invoice
	input := map[string]interface{}{
		"tenantid":       	_mastende.Tenant.ID,
		"date":       		"2017-08-01",
		"duedate": 		"2017-08-05",
	}

	if err := _mastende.CreateInvoice(input); err != nil {

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

func testCreatePeriodRange(t *testing.T)  {

	if err := periods.CreateFinancialPeriodRange(&a, "2017-07-01", 12); err != nil{

		t.Error("Error creating period data!")
		return
	}
}
