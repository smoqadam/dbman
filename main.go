package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	app := NewApp()
	r := mux.NewRouter()

	dir := "./static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/", app.home)
	r.HandleFunc("/login", app.login).Methods("POST")

	auth := r.NewRoute().Subrouter()
	auth.Use(app.auth)
	auth.Use(app.conn)
	auth.HandleFunc("/databases", app.dblist)
	auth.HandleFunc("/db/{db}", app.tablelist)
	auth.HandleFunc("/table/{db}/{table}", app.table)
	auth.HandleFunc("/table/{db}/{table}/data", app.data)

	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	// r.Handle("/assets/css", http.FileServer(http.Dir("./assets/css")))
	// // // r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Fatal(http.ListenAndServe(":3001", r))
}
