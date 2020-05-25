package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/sessions"
)

var (
	cookieName   = "user"
	cookieSecret = []byte("secret")
)

type App struct {
	db     Database
	Conn   *sql.DB
	cookie *sessions.CookieStore
}

func NewApp() *App {
	return &App{
		cookie: sessions.NewCookieStore(cookieSecret),
	}
}

func (app *App) json(w http.ResponseWriter, out interface{}) {
	w.Header().Add("Content-Type", "application/json")
	res, _ := json.Marshal(out)
	w.Write(res)
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	// Add whatever other types you need
	default:
		return ""
	}
}

func (app *App) html(w http.ResponseWriter, file string, data interface{}) {
	base := path.Join("templates", "base.html")
	header := path.Join("templates", "header.html")
	footer := path.Join("templates", "footer.html")
	view := path.Join("templates", file)

	tpl, err := template.New("").Funcs(template.FuncMap{"ToString": ToString}).ParseFiles(base, header, footer, view)
	if err != nil {
		data = struct{ Msg string }{
			Msg: fmt.Sprintf("%s: %s", "View not found", view),
		}
		view = path.Join("templates", "error.html")
		tpl, _ = template.ParseFiles(base, header, footer, view)

	}
	tpl.ExecuteTemplate(w, "base", data)
}

func (app *App) connect(s string, h string, u string, p string) error {
	switch s {
	case "mysql":
		conn := fmt.Sprintf("%s:%s@tcp(%s)/?", u, p, h)
		fmt.Println(conn)
		db, err := sql.Open("mysql", conn)
		if err != nil {
			return err
		}
		d := &Mysql{}

		app.db = d
		app.Conn = db
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
		log.Println("err", err)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
