package helper

import (
	"fmt"
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//insert into table
func InsertIntoTable(conn *pgxpool.Pool, table string, columns []string, values []interface{}) (string, error) {
	len, len2 := len(columns), len(values)
	if len != len2 {
		return "", fmt.Errorf("number of columns and values must be the same")
	}
	cols := joinColumns(columns)
	vals := joinValues(values)
	//fmt.Printf("Inserting into table %s with columns %s and values %s\n", table, cols, vals)
	//fmt.Println("Query is ", query)
	tag := conn.QueryRow(context.Background(), `INSERT INTO $1 ($2) VALUES ($3);`, table, cols, vals)
	var row string
	err := tag.Scan(&row)
	if err != nil {
		return "", err
	}
	return row, nil
}

func joinColumns(columns []string) string {
	return strings.Join(columns, ", ")
}

func joinValues(values []interface{}) string {
	var strValues []string
	for _, val := range values {
		switch v:= val.(type) {
		case string:
			strValues = append(strValues, fmt.Sprintf("'%s'", v))
		default:
			strValues = append(strValues, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(strValues, ", ")
}
