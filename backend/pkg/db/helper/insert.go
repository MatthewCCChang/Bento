package helper

import (
	"fmt"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//insert into table
func InsertIntoTable(conn *pgxpool.Pool, table string, columns []string, values []interface{}, retVals []string, ret bool) (string, error) {
	len, len2 := len(columns), len(values)
	if len != len2 {
		return "", fmt.Errorf("number of columns and values must be the same")
	}
	cleanTable := pgx.Identifier{table}.Sanitize()
	fmt.Println(cleanTable)
	cols := JoinColumns(columns)
	vals := JoinValues(values)
	rets := JoinColumns(retVals)
	//fmt.Printf("Inserting into table %s with columns %s and values %s\n", table, cols, vals)
	//fmt.Println("Query is ", query)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, cleanTable, cols, vals)
	//if return val requested
	if ret{
		query += fmt.Sprintf(`RETURNING %s`, rets)
	}
	query += `;`
	tag := conn.QueryRow(context.Background(), query)
	var row string
	err := tag.Scan(&row)
	if err != nil {
		return "", err
	}
	return row, nil
}

func JoinColumns(columns []string) string {
	var res []string
	for _, col := range columns{
		clean := pgx.Identifier{col}.Sanitize()
		res = append(res, clean)
	}
	return strings.Join(res, ", ")
}

func JoinValues(values []interface{}) string {
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
