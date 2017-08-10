package periods_test

import (
	"github.com/khosimorafo/mastende/db"
	"testing"
	"github.com/jinzhu/configor"
	"os"
	"github.com/khosimorafo/mastende/periods"
)

var a db.App

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

func TestPeriodRead(t *testing.T) {

	// Pass collection to package
	testCreatePeriodRange(t)

	period, _ := periods.NewInstanceWithDate(&a, "2017-08-10")

	t.Log(period)
}

func testCreatePeriodRange(t *testing.T)  {

	if err := periods.CreateFinancialPeriodRange(&a, "2017-07-01", 12); err != nil{

		t.Error("Error creating period data!")
		return
	}
}