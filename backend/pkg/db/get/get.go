package get

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JSONB []byte
type AggregatedItem struct {
	Name string //item + modifier
	Price float32
	Quantity int
	UserIds []int
	Session_id int
}

//GetSessionOrder retrieves all rows belonging to the same order
func GetSessionOrder(conn *pgxpool.Pool, sessionId int) (map[string]interface{}, error) {
	query := fmt.Sprintf(`SELECT * AS total FROM orders WHERE id = %d;`, sessionId)
	rows, err := conn.Query(context.Background(), query)

	res := make(map[string]interface{}) 
	defer rows.Close()

	//iterate through rows and aggregate items with same name and modifiers, also keep track of quantity and user ids for each item
	for rows.Next() {
		var id int
		var session_id int
		var user_id int
		var item string
		var modifiers JSONB
		var price float32

		err = rows.Scan(&id, &session_id, &user_id, &item, &modifiers, &price)
		if err != nil {
			return nil, err
		}
		key := item + string(modifiers) //combine item and modifiers to create unique key for aggregation
		if _, ok := res[key]; !ok {
			res[key] = AggregatedItem{
				Name: key,
				Price: price,
				Quantity: 1,
				UserIds: []int{user_id},
				Session_id: session_id,
			}
		}else{
			aggItem := res[key].(AggregatedItem)
			aggItem.Quantity += 1
			aggItem.UserIds = append(aggItem.UserIds, user_id)
			res[key] = aggItem
		}

	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	
	
	return res, nil
}