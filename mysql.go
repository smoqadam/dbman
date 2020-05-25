package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct{}

func (m Mysql) DbList(db *sql.DB) []string {
	dblist := []string{}
	rows, err := db.Query("Show databases")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var dbName []uint8
		if err := rows.Scan(&dbName); err != nil {
			panic(err)
		}
		dblist = append(dblist, string(dbName))
	}
	return dblist
}

func (m Mysql) TableList(dbName string, db *sql.DB) []string {
	dblist := []string{}
	_, err := db.Exec("USE " + dbName)

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("Show tables;")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var dbName []uint8
		if err := rows.Scan(&dbName); err != nil {
			panic(err)
		}
		dblist = append(dblist, string(dbName))
	}
	return dblist
}

type Describe struct {
	Field   []byte
	Type    []byte
	Null    []byte
	Key     []byte
	Default []byte
	Extra   []byte
}

func (m Mysql) Describe(dbName string, table string, db *sql.DB) []Describe {

	rows := []Describe{}

	_, err := db.Exec("Use " + dbName)

	if err != nil {
		panic(err)
	}

	r, err := db.Query("DESCRIBE " + table)
	if err != nil {
		panic(err)
	}
	// rd := []Describe{}
	for r.Next() {
		// row := make([]interface{}, 6)
		rdr := Describe{}
		if err := r.Scan(&rdr.Field, &rdr.Type, &rdr.Null, &rdr.Key, &rdr.Default, &rdr.Extra); err != nil {
			panic(err)
		}
		fmt.Println(rdr)
		rows = append(rows, rdr)
	}
	return rows
}
