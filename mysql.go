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

// Tables return an array of table of the given db
func (m *Mysql) Tables(dbName string) ([]Table, error) {
	dblist := []Table{}

	rows, err := m.Query("information_schema", "SELECT TABLE_NAME, ENGINE FROM information_schema.tables WHERE TABLE_SCHEMA = ?", dbName)
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

	rows, err := m.Query(dbName, "DESCRIBE "+table)
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

func (m *Mysql) Query(dbName string, query string, args ...interface{}) (*sql.Rows, error) {

	fmt.Println("query", query, args)
	_, err := m.DB().Exec("Use " + dbName)
	if err != nil {
		return nil, err
	}
	stmt, err := m.DB().Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Query(args...)
}

func (m *Mysql) Data(dbName string, query string, args ...interface{}) (Rows, error) {
	rows := Rows{}

	r, err := m.Query(dbName, query, args...)
	if err != nil {
		return rows, err
	}

	cols, _ := r.Columns()
	rows.Fields = cols
	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
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
		row := rw
		rows.Values = append(rows.Values, row)
	}
	return rows, nil
}
