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
	r.HandleFunc("/", app.home)
	r.HandleFunc("/login", app.login).Methods("POST")

	auth := r.NewRoute().Subrouter()
	auth.Use(app.auth)
	auth.Use(app.conn)
	auth.HandleFunc("/databases", app.dblist)
	auth.HandleFunc("/db/{db}", app.tablelist)
	auth.HandleFunc("/table/{db}/{table}", app.table)

	log.Fatal(http.ListenAndServe(":3001", r))
}
