package items

import (
	"github.com/khosimorafo/mastende/db"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/khosimorafo/mastende/utils"
	"log"
	"fmt"
)

var a db.App

/**
*
* Create new items record.
*/
func New(app *db.App) *Item {

	a = *app
	a.SetCollection("items")
	t := &Item{}
	return t
}

/***
*
* When provided an id, This function returns a items with data already read from the database
*
*/
func NewInstanceWithId(app *db.App, id string) (*Item, error) {

	a = *app
	item := New(&a)
	item.ID = id

	if err := item.Read(); err != nil {

		return item, errors.New("Error while attempting to read items.")

	}
	return item, nil
}

/**
*
* Create item list.
*/
func ItemList(app *db.App, t *[]Item) error {

	a = *app
	a.SetCollection("items")

	if err := GetItemList(t); err !=nil{

		return errors.New("Error while attempting to create item list.")
	}

	fmt.Printf("\nGetItemList : ", &t)

	return nil
}

/**
*
* Create item list by name.
*/
func ItemListByName(app *db.App, t *[]Item, search_str string) error {

	a = *app
	a.SetCollection("items")

	if err := GetItemListByName(t, search_str); err !=nil{

		return errors.New("Error while attempting to create item list.")
	}

	fmt.Printf("\nGetItemList : ", &t)

	return nil
}

func (item *Item) Persist() (error) {

	item.ID = utils.RandStringRunes(25)

	if err := persist(item); err != nil{

		return errors.New("Error while attempting to save items.")
	}
	return nil
}

func (item *Item) Read() (error) {

	if err := get(item); err != nil{

		return errors.New("Error retrieving items.")
	}
	return nil
}

func (item *Item) Delete() (error) {

	if err := delete(item); err != nil{

		return errors.New("Error retrieving items.")
	}
	return nil
}

func (item *Item) Update() (error) {

	if err := update(item); err != nil{

		return errors.New("Error updating items.")
	}
	return nil
}

func (item *Item) Validate() (error) {


	return nil
}

func (item *Item) SetStatusActive() (error) {

	// Setup update
	colQuerier := bson.M{"id": item.ID}
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

func (item *Item) SetStatusInActive() (error) {

	// Setup update
	colQuerier := bson.M{"id": item.ID}
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

func persist(item *Item) (error) {

	item.UpdatedTime = time.Now().String()
	item.CreatedTime = time.Now().String()

	fmt.Printf("\nItem to be saved : ", item)

	if err := a.Collection.Insert(item); err != nil {
		log.Println(err)
		return errors.New("Duplicate ZAID submitted.")
	}
	return nil
}

func get(item *Item) error {

	if err := a.Collection.Find(bson.M{"id": item.ID}).One(&item); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func delete(item *Item) error {

	if err := a.Collection.Remove(bson.M{"id": item.ID}); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func update(item *Item) error{

	// Setup update
	colQuerier := bson.M{"id": item.ID}
		change := bson.M{"$set": bson.M{
			"name": 	item.Name,
			"description":	item.Description,
			"rate":		item.Rate,
			"status":	item.Status,
			"updatedtime": 	time.Now().String(),
		},
	}

	item.UpdatedTime = time.Now().String()

	if err := a.Collection.Update(colQuerier, change); err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}

	return nil
}

func GetItemList (items *[]Item) (error) {

	list := []Item{}

	if err := a.Collection.Find(bson.M{}).All(&list); err != nil {

		return err

	}
	//Reassign the address pointed to by the items and return
	*items = list; return nil
}

func GetItemListByName (items *[]Item, search_str string) (error) {

	list := []Item{}

	if err := a.Collection.Find(bson.M{"name": search_str}).All(&list); err != nil {

		return err
	}
	//Reassign the address pointed to by the items and return
	*items = list; return nil
}

type Item struct {

	ID          	string  	`json:"id",omitempty"`
	Name        	string  	`json:"name"`
	Description 	string  	`json:"description                                                                  "`
	Rate        	float64 	`json:"rate"`
	Status        	string 		`json:"status"`

	CreatedTime 	string  	`json:"createdtime"`
	UpdatedTime 	string  	`json:"updatedtime"`
}
