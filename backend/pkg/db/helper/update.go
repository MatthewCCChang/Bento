package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//insert into table
func UpdateTable(conn *pgxpool.Pool, table string, columns, cond []string, values, condVals []interface{}, retVals []string, ret bool) (map[string]interface{}, error) {
	length, length2 := len(columns), len(values)
	if length != length2 {
		return map[string]interface{}{}, fmt.Errorf("number of columns and values must be the same")
	}
	cleanTable := pgx.Identifier{table}.Sanitize()
	fmt.Println(cleanTable)
	rets := JoinColumns(retVals)

	var query strings.Builder
	query.Grow(256)

	query.WriteString(fmt.Sprintf(`UPDATE %s SET `, cleanTable))

	//store key=value statements with parameterized values
	var pairs []string
	var args []any
	
	for i, col := range columns {
        // ALWAYS sanitize column names too!
        cleanCol := pgx.Identifier{col}.Sanitize() 
        
        // i is 0-indexed, but SQL parameters are 1-indexed
        pairs = append(pairs, fmt.Sprintf(`%s=$%d`, cleanCol, i+1))
        args = append(args, values[i])
    }
	//join parmaterized statements
	query.WriteString(strings.Join(pairs, `, `))
	colCount := len(columns)

	query.WriteString(fmt.Sprintf(` WHERE `))

	var condition []string
	//add WHERE clause can have multiple
	for i, col := range cond {
		cleanCol := pgx.Identifier{col}.Sanitize()

		condition = append(condition, fmt.Sprintf(`%s=$%d`, cleanCol, colCount+i+1))
        args = append(args, condVals[i])
	}
	query.WriteString(strings.Join(condition, ` AND `))
	

	//if return val requested
	if ret{
		query.WriteString(fmt.Sprintf(` RETURNING %s`, rets))
	}
	query.WriteByte(';')

	//execute query with arguments as $1, $2...etc
	rows, err := conn.Query(context.Background(), query.String(), args...)
	if err != nil {
		return map[string]interface{}{}, err
	}

	res, err := pgx.RowToMap(rows)
	if err != nil{
		return map[string]interface{}{}, err
	}

	return res, nil
}


func JoinColumnsWithValues(columns []string, values []interface{}) (map[string]any, error) {  
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
