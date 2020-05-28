package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"test": "test",
	}
	app.Response.html(w, "views/login.html", data)
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	host := r.Form.Get("host")
	username := r.Form.Get("username")
	pass := r.Form.Get("password")
	err := app.connect("mysql", host, username, pass)
	if err != nil {
		//TODO: Change panic to redirect to login
		panic(err)
	}

	session, _ := app.cookie.Get(r, cookieName)
	session.Values["authenticated"] = true
	session.Values["host"] = host
	session.Values["user"] = username
	session.Values["pass"] = pass
	session.Save(r, w)

	http.Redirect(w, r, "/databases", 302)
}

func (app *App) dblist(w http.ResponseWriter, r *http.Request) {
	log.Println("DB", app.db.DB())
	log.Println("DB", app.Conn)
	dbs, err := app.db.Databases()
	if err != nil {
		panic(err)
	}
	app.Response.html(w, "views/databases.html", dbs)
}

func (app *App) tablelist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tables, err := app.db.Tables(vars["db"])
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"Db":     vars["db"],
		"Tables": tables,
	}
	app.Response.html(w, "views/table/list.html", data)
}

func (app *App) table(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Println("TABLE")
	cols, err := app.db.Columns(vars["db"], vars["table"])
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"DbName":  vars["db"],
		"Table":   vars["table"],
		"Columns": cols,
	}
	app.Response.html(w, "views/table/structure.html", data)
}

func (app *App) data(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	rows, err := app.db.Data(vars["db"], "select * from "+vars["table"])
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"DbName": vars["db"],
		"Table":  vars["table"],
		"rows":   rows,
	}
	app.Response.html(w, "views/table/data.html", data)
}
