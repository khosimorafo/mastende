package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
	"errors"
)

var a App

type App struct {

	Session *mgo.Session
	Collection *mgo.Collection
	Database *mgo.Database
}

func DB(uri string, db string ) *App {

	a.Session = AppCollection(uri)
	a.Database = a.Session.DB(db)

	return &a
}

func (a *App)SetCollection(coll string){

	a.Collection = a.Database.C(coll)
}

func AppCollection(uri string) (*mgo.Session) {

	//uri := "mongodb://mastende:mastende@ds115573.mlab.com:15573/mastende-test"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	//defer sess.Close()

	//sess.SetSafe(&mgo.Safe{})

	return sess;
}

func Update(db *App, sel interface{}, change interface{}) error {

	if err := db.Collection.Update(sel, change); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

