package items_test

import (
	"os"
	"testing"
	"fmt"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/db"
	"github.com/jinzhu/configor"
)

var a db.App

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {

	// Init test session/db/collection
	configor.Load(&db.Config, "config.yml")
	a = *db.DB(db.Config.DB.DbUrl, db.Config.DB.DbName)
	a.Database.C("items").DropCollection()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	a.Session.Close()

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestItemCreate(t *testing.T) {

	// Pass collection to package

	itemID := testCreate(t)
	testGet(itemID, t)
	testUpdate(itemID, t)
	testGet(itemID, t)
	testDelete(itemID, t)
}

func testCreate(t *testing.T) string {

	item := items.New(&a)
	item.Name = "Test Item Name"
	item.Rate = 330

	if err := item.Persist(); err != nil {
		t.Errorf("Error creating items", err.Error())
	}

	return item.ID
}

func testGet(itemId string, t *testing.T) {

	item, err := items.NewInstanceWithId(&a, itemId)
	if err != nil {
		t.Errorf("Error retrieving items : ", err.Error())
	}

	fmt.Printf("Retrieved items is : ", item)
}

func testUpdate(itemId string, t *testing.T) {

	item, err := items.NewInstanceWithId(&a, itemId)
	if err != nil {
		t.Errorf("Error retrieving items : ", err.Error())
	}

	update_name := "Updated Item Name"
	item.Name = update_name
	item.Rate = 333.0
	item.Update()

	item.Read()

	if item.Name != update_name {
		t.Errorf("Expected the name to change to (%v). Got %v", update_name, item.Name)
	}

}

func testDelete(itemId string, t *testing.T) {

	item, err := items.NewInstanceWithId(&a, itemId)
	if err != nil {
		t.Errorf("Error getting items for delete : ", err.Error())
	}

	if err = item.Delete(); err != nil {
		t.Errorf("Error deleting items", err.Error())
		return
	}
}

func TestListItemsByName(t *testing.T) {

	testCreate(t)

	var i []items.Item

	err := items.ItemListByName(&a, &i, "Test Item Name")

	if err != nil {

		t.Error("Error reading items")
	}

	t.Log(i)
}


