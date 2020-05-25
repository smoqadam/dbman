package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	db *sql.DB
}

func (m *Mysql) SetDB(db *sql.DB) {
	m.db = db
}

func (m *Mysql) DB() *sql.DB {
	return m.db
}

func (m *Mysql) Databases() ([]string, error) {
	dblist := []string{}
	rows, err := m.DB().Query("Show databases")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dbName []uint8
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		dblist = append(dblist, string(dbName))
	}
	return dblist, nil
}

// Tables return an array of table
func (m *Mysql) Tables(dbName string) ([]Table, error) {
	dblist := []Table{}
	_, err := m.DB().Exec("USE  information_schema")

	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT TABLE_NAME, ENGINE FROM information_schema.tables WHERE TABLE_SCHEMA =  '%s'", dbName)
	rows, err := m.DB().Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Name, &table.Engine); err != nil {
			return nil, err
		}
		dblist = append(dblist, table)
	}
	return dblist, nil
}

func (m *Mysql) Columns(dbName string, table string) ([]Column, error) {

	columns := []Column{}

	_, err := m.DB().Exec("Use " + dbName)

	if err != nil {
		return nil, err
	}

	rows, err := m.DB().Query("DESCRIBE " + table)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		col := Column{}
		if err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra); err != nil {
			return nil, err
		}
		fmt.Println(col)
		columns = append(columns, col)
	}
	return columns, nil
}

func (m *Mysql) Query(dbName string, query string) (Rows, error) {
	// var rows []map[string]interface{}
	rows := Rows{}
	_, err := m.DB().Exec("Use " + dbName)

	if err != nil {
		return rows, err
	}

	r, err := m.DB().Query(query)
	if err != nil {
		panic(err)
	}

	cols, _ := r.Columns()
	rows.Fields = cols
	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for r.Next() {
		r.Scan(dest...)
		rw := make([]string, len(cols))
		for i, raw := range rawResult {

			if raw == nil {
				rw[i] = "null"
			} else {
				rw[i] = string(raw)
			}
		}
		row := Row{
			Values: rw,
		}
		rows.Values = append(rows.Values, row)
		fmt.Printf("ROW: %#v\n", rows)
	}
	return rows, nil
}
