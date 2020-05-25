package main

import (
	"database/sql"
)

type Database interface {
	DbList(*sql.DB) []string
	TableList(string, *sql.DB) []string
	Describe(string, string, *sql.DB) []Describe
}
