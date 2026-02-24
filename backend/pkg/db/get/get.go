package get

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type JSONB []byte
type Item struct{
	Id int
	Version_id int
	Name string
	Description string
	Price float32
	Category string
	Modifiers JSONB
}
type AggregatedItem struct {
	Name string //item + modifier
	Price float32
	Quantity int
	UserIds []int
	Session_id int
}

//GetSessionOrder retrieves all rows belonging to the same order
func GetSessionOrder(conn *pgxpool.Pool, sessionId int) (map[string]interface{}, error) {
	rows, err := conn.Query(context.Background(), `SELECT * AS total FROM orders WHERE id =$1;`, sessionId)

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
		key := item + " " + string(modifiers) //combine item and modifiers to create unique key for aggregation
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


//majorty will be redis
//GetMenu - get from items table where version id is equal to blah if not cached already
func GetMenu(ctx context.Context, rdb *redis.Client, conn *pgxpool.Pool, restaurant_id int) ([]Item, error){
	//get curr version for rest
	var version string
	//chekc if version exists in redis alr  
	version, err := rdb.Get(ctx, fmt.Sprintf("restaurant:%d:active", restaurant_id)).Result()  //get the curr active
	//if yes, return the cached json
	//if not, get menu id and version id then fetch items
	if err != nil{
		//fetch from db 
		row := conn.QueryRow(context.Background(), `SELECT id FROM version as v LEFT JOIN menu as m ON m.id=v.menu_id LEFT JOIN restaurant as r ON r.id=m.restaurant_id WHERE r.id=$1 AND v.is_active=true;`, restaurant_id)
		err := row.Scan(&version)
		if err != nil{
			return []Item{}, fmt.Errorf("Error retreiving version %w", err)
		}
		//update redis with new one as well
	}
	
	json, err := rdb.Get(ctx, fmt.Sprintf("restaurant:%d:menu:v%s", restaurant_id, version)).Result()
	if err != nil{
		return []Item{}, fmt.Errorf("Error retrieiving menu items %w", err)
	}
	fmt.Printf("%s", json)
	
	//turn json into Items

	return []Item{}, nil
}