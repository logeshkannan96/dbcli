package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}

func IsConnected() bool {
	if db == nil {
		return false
	}
	return db.Ping() == nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func ExecuteQuery(query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			var v interface{}
			values[i] = &v
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		// Create a map for the row data
		row := make(map[string]interface{})
		for i, column := range columns {
			val := *(values[i].(*interface{}))
			row[column] = val
		}
		results = append(results, row)
	}
	return results, nil
}
