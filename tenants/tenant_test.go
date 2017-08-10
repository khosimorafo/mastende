package tenants_test

import (
	"os"
	"testing"
	"fmt"
	"github.com/khosimorafo/mastende/tenants"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
	"github.com/khosimorafo/mastende/utils"
)

var a db.App

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {

	// Init test session/db/collection
	configor.Load(&db.Config, "config.yml")
	a = *db.DB(db.Config.DB.DbUrl, db.Config.DB.DbName)
	a.Database.C("tenants").DropCollection()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	a.Session.Close()

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestTenantCreate(t *testing.T) {

	// Pass collection to package

	tenantID := testCreate(t)
	testGet(tenantID, t)
	testUpdate(tenantID, t)
	testDelete(tenantID, t)
	//testGet(tenantID, t)
}

func testCreate(t *testing.T) string {

	tenant := tenants.New(&a)
	tenant.Name = "Test Tenant Name Name"
	tenant.ZAID = utils.RandNumberRunes(13)

	if err := tenant.Persist(); err != nil {
		t.Errorf("Error creating tenants", err.Error())
	}

	return tenant.ID
}

func testGet(tenantId string, t *testing.T) {

	tenant, err := tenants.NewInstanceWithId(&a, tenantId)
	if err != nil {
		t.Errorf("Error retrieving tenants : ", err.Error())
	}

	fmt.Printf("Retrieved tenants is : ", tenant)
}

func testUpdate(tenantId string, t *testing.T) {

	tenant, err := tenants.NewInstanceWithId(&a, tenantId)
	if err != nil {
		t.Errorf("Error retrieving tenants : ", err.Error())
	}

	update_name := "Updated Name"
	tenant.Name = update_name
	tenant.Update()

	tenant.Read()

	if tenant.Name != update_name {
		t.Errorf("Expected the name to change to (%v). Got %v", update_name, tenant.Name)
	}

}

func testDelete(tenantId string, t *testing.T) {

	tenant, err := tenants.NewInstanceWithId(&a, tenantId)
	if err != nil {
		t.Errorf("Error getting tenants for delete : ", err.Error())
	}

	if err = tenant.Delete(); err != nil {
		t.Errorf("Error deleting tenants", err.Error())
		return
	}
}

