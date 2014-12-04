package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

type SalesRep struct {
	Name      string `form:"name"`
	BuzzWords string `form:"description"`
}

// DB Returns a martini.Handler
func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("SalesRepBingoData"))
		defer s.Close()
		c.Next()
	}
}

// GetAll returns all Wishes in the database
func GetAll(db *mgo.Database) []SalesRep {
	var reps []SalesRep
	db.C("salesReps").Find(nil).All(&reps)
	return reps
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	//START3 OMIT
	m.Use(DB())

	m.Get("/salesreps", func(r render.Render, db *mgo.Database) {
		r.JSON(200, GetAll(db))
	})

	m.Post("/salesreps", binding.Form(SalesRep{}), func(rep SalesRep, r render.Render, db *mgo.Database) {
		db.C("salesReps").Insert(rep)
		r.HTML(200, "list", GetAll(db))
	})

	m.Run()
}
