package main

import (
	"database/sql"
)

// Database interface is resposible to get information about
// the database engine such as Mysql, sqlite, etc.
type Database interface {
	SetDB(*sql.DB)

	DB() *sql.DB
	// Return list of databases
	Databases() ([]string, error)

	// Tables receives a db nam and  then
	// return an array of Table
	Tables(string) ([]Table, error)

	// Columns receives db name, table name and a pointer to
	// sql.DB then return an array of Columns
	Columns(string, string) ([]Column, error)

	Query(string, string) (Rows, error)
}

type Column struct {
	Field   []byte
	Type    []byte
	Null    []byte
	Key     []byte
	Default []byte
	Extra   []byte
}

type Table struct {
	Name    string
	Engine  string
	Columns []Column
}

type Rows struct {
	Fields []string
	Values []Row
}

type Row struct {
	Values []string
}
