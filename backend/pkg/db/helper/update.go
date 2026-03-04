package helper

import (
	"fmt"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//insert into table
func UpdateTable(conn *pgxpool.Pool, table string, columns []string, values []interface{}, retVals []string, ret bool) (string, error) {
	len, len2 := len(columns), len(values)
	if len != len2 {
		return "", fmt.Errorf("number of columns and values must be the same")
	}
	cleanTable := pgx.Identifier{table}.Sanitize()
	fmt.Println(cleanTable)
	res, err := joinColumnsWithValues(columns, values)
	if err != nil{
		return "", err
	}
	rets := joinColumns(retVals)

	var query strings.Builder
	query.Grow(256)

	query.WriteString(fmt.Sprintf(`UPDATE %s SET`, cleanTable))

	//store key=value statements with parameterized values
	var pairs []string
	var args []any
	i := 1
	//dynamically add col-value pairs
	for col, val := range res{
		pairs = append(pairs, fmt.Sprintf(`%s=$%d`, col, i))
		args = append(args, val)
		i++
	}

	//join parmaterized statements
	query.WriteString(strings.Join(pairs, `, `))
	//add WHERE clause
	query.WriteString(fmt.Sprintf(` WHERE id=$%d`, i))

	//if return val requested
	if ret{
		query.WriteString(fmt.Sprintf(` RETURNING %s`, rets))
	}
	query.WriteByte(';')

	//execute query with arguments as $1, $2...etc
	tag := conn.QueryRow(context.Background(), query.String(), args...)
	var row string
	err = tag.Scan(&row)
	if err != nil {
		return "", err
	}
	return row, nil
}


func joinColumnsWithValues(columns []string, values []interface{}) (map[string]any, error) {  
	if len(columns) != len(values){
		return map[string]any{}, fmt.Errorf("Error: length of column and values don't match")
	}

	//value can be of multiple types
	res := make(map[string]any)
	for idx, val := range values {
		col := columns[idx]
		clean := pgx.Identifier{col}.Sanitize()
		var value string
		switch v:= val.(type) {
		case string:
			value =  fmt.Sprintf("'%s'", v)
		default:
			value = fmt.Sprintf("%v", v)
		}
		
		//store key-value pair as col-value
		res[clean] = value 
	}
	return res, nil
}
