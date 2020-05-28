package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	cookieName   = "user"
	cookieSecret = []byte("secret")
	dbp          *sql.DB
)

type App struct {
	db       Database
	Conn     *sql.DB
	cookie   *sessions.CookieStore
	Response *Response
}

func NewApp() *App {
	return &App{
		Response: &Response{},
		cookie:   sessions.NewCookieStore(cookieSecret),
	}
}

func (app *App) connect(s string, h string, u string, p string) error {
	switch s {
	case "mysql":
		conn := fmt.Sprintf("%s:%s@tcp(%s)/?", u, p, h)
		db, err := sql.Open("mysql", conn)
		if err != nil {
			return err
		}
		database := &Mysql{}
		database.SetDB(db)
		app.Conn = db
		app.db = database
		dbp = db
	}
	return nil
}

func (app *App) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.cookie.Get(r, cookieName)
		log.Println("AUTH")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {

			next.ServeHTTP(w, r)
		}
	})
}

func (app *App) conn(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CONN")
		session, _ := app.cookie.Get(r, cookieName)
		host := session.Values["host"]
		user := session.Values["user"]
		pass := session.Values["pass"]

		err := app.connect("mysql", host.(string), user.(string), pass.(string))
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
