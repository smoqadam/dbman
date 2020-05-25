package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	app.html(w, "views/login.html", nil)
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
	app.html(w, "views/databases.html", app.db.DbList(app.Conn))
}

func (app *App) tablelist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("db: ", vars["db"])
	data := map[string]interface{}{
		"Db":     vars["db"],
		"Tables": app.db.TableList(vars["db"], app.Conn),
	}
	fmt.Println(data)
	app.html(w, "views/table/list.html", data)
}

func (app *App) table(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := map[string]interface{}{
		"DbName": vars["db"],
		"Table":  vars["table"],
		"Rows":   app.db.Describe(vars["db"], vars["table"], app.Conn),
	}
	app.html(w, "views/table/structure.html", data)
}
