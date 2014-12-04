package main

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name  string
	Phone string
}

var (
	mgoSession   *mgo.Session
	databaseName = "myDB"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("localhost")
		if err != nil {
			panic(err) //no, not really
		}
	}
	return mgoSession.Clone()
}

func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(databaseName).C(collection)
	return s(c)
}

func main() {
	m := martini.Classic()
	m.Get("/:name", func(params martini.Params) string {
		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()

		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)

		c := session.DB("test").C("people")
		result := Person{}
		err = c.Find(bson.M{"name": params["name"]}).One(&result)
		if err != nil {
			return "whoops"
		}

		return "Phone:" + result.Phone
	})
	m.Post("/:name", func(params martini.Params) string {
		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		c := session.DB("test").C("people")
		err = c.Insert(&Person{"Alisha", "+55 53 8116 9639"},
			&Person{"Craig", "+55 53 8402 8510"})
		if err != nil {
			log.Fatal(err)
		}
		return "Success!"
	})
	m.Run()
}

func SetupDB() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("localhost")
		if err != nil {
			PanicIf(err)
		}
	}
	return mgoSession.Clone()
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
