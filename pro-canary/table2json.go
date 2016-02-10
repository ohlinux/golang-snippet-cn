package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

//所有的都变成string了,对值使用relect没有作用
func dumpTable(table string) {
	// ...''
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "12345", "127.0.0.1", "3306", "orp")
	db, _ := sql.Open("mysql", conn)
	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	//fmt.Println(err)
	columns, _ := rows.Columns()
	//fmt.Println(err)

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(scanArgs...)

		record := make(map[string]interface{})

		for i, col := range values {
			if col != nil {
				record[columns[i]] = fmt.Sprintf("%s", string(col.([]byte)))
			}
		}

		s, _ := json.Marshal(record)
		os.Stdout.Write(s)
	}
}

func main() {
	dumpTable("module")
}
