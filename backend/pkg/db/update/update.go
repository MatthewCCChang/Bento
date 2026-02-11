package update

import (
	"fmt"
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//insert into table
func InsertIntoTable(conn *pgxpool.Pool, table string, columns []string, values []interface{}) (int64, error) {
	len, len2 := len(columns), len(values)
	if len != len2 {
		return 0, fmt.Errorf("number of columns and values must be the same")
	}
	cols := joinColumns(columns)
	vals := joinValues(values)
	//fmt.Printf("Inserting into table %s with columns %s and values %s\n", table, cols, vals)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s);`, table, cols, vals)
	//fmt.Println("Query is ", query)
	tag, err := conn.Exec(context.Background(), query)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
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