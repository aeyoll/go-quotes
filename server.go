package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/gzip"
  "labix.org/v2/mgo"
)

type Quote struct {
  Content   string  `form:"content"`
}

func DB() martini.Handler {
  session, err := mgo.Dial("mongodb://localhost")
  if err != nil {
    panic(err)
  }

  return func(c martini.Context) {
    s := session.Clone()
    c.Map(s.DB("quotes"))
  }
}

func GetAll(db *mgo.Database) []Quote {
  var quotes []Quote
  db.C("quotes").Find(nil).All(&quotes)
  return quotes
}

func main() {
  m := martini.Classic()

  m.Use(gzip.All())

  // render html templates from templates directory
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))

  m.Use(DB())

  m.Get("/", func(r render.Render, db *mgo.Database) {
    data := map[string]interface{}{"quotes": GetAll(db)}
    r.HTML(200, "list", data)
  })

  m.Post("/", binding.Form(Quote{}), func(r render.Render, db *mgo.Database, quote Quote) {
    db.C("quotes").Insert(quote)
    data := map[string]interface{}{"quotes": GetAll(db)}
    r.HTML(200, "list", data)
  })

  m.Run()
}
