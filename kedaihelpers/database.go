package kedaihelpers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DBStruct adalah struct utk interface sqlx.DB
type DBStruct struct {
	Dbx *sqlx.DB
}

func mapBytesToString(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
	return m
}

// DatabaseQueryRows berfungsi mendapatkan hasil query utk yg lebih dari 1 baris
func (dbs *DBStruct) DatabaseQueryRows(query string, args ...interface{}) []map[string]interface{} {

	var datarows []map[string]interface{}
	rows, err := dbs.Dbx.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			if err != nil {
				fmt.Println("Scan Error", err)
			}
			datarows = append(datarows, mapBytesToString(results))
		}
	}

	return datarows
}

// DatabaseQuerySingleRow berfungsi mendapatkan hasil query utk yg hanya 1 baris
func (dbs *DBStruct) DatabaseQuerySingleRow(query string, args ...interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	rows, err := dbs.Dbx.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			if err != nil {
				fmt.Println("Scan Error", err)
			}
			return mapBytesToString(results)
		}
	}

	return result
}
