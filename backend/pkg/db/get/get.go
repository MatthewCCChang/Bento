package get

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/MatthewCCChang/Bento/backend/pkg/db/helper"
)

type Item struct{
	Id int	`db:"id"`
	Version_id int	`db:"version_id"`
	Name string	`db:"name"`
	Description string	`db:"description"`
	Price float64	`db:"price"`
	Category string	`db:"category"`
	Modifiers json.RawMessage	`db:"modifiers"`
}
type AggregatedItem struct {
	Name string //item + modifier
	Price float32
	Quantity int
	UserIds []int
	Session_id int
}

func GetRow(conn *pgxpool.Pool, table string, cond []string, vals []interface{}, retVals []string) (map[string]interface{}, error){
	cleanTable := pgx.Identifier{table}.Sanitize()

	rets := helper.JoinColumns(retVals)

	var query strings.Builder
	query.Grow(256)

	query.WriteString(fmt.Sprintf(`SELECT %s FROM %s WHERE `, rets, cleanTable))

	//store key=value statements with parameterized values
	var pairs []string
	var args []any
	
	for i, col := range cond {
        // ALWAYS sanitize column names too!
        cleanCol := pgx.Identifier{col}.Sanitize() 
        
        // i is 0-indexed, but SQL parameters are 1-indexed
        pairs = append(pairs, fmt.Sprintf(`%s=$%d`, cleanCol, i+1))
        args = append(args, vals[i])
    }

	//join parmaterized statements
	query.WriteString(strings.Join(pairs, `, `))
	query.WriteByte(';')

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


//GetSessionOrder retrieves all rows belonging to the same order
func GetSessionOrder(conn *pgxpool.Pool, sessionId int) (map[string]interface{}, error) {
	rows, err := conn.Query(context.Background(), `SELECT * FROM orders WHERE id =$1;`, sessionId)

	res := make(map[string]interface{}) 
	defer rows.Close()

	//iterate through rows and aggregate items with same name and modifiers, also keep track of quantity and user ids for each item
	for rows.Next() {
		var id int
		var session_id int
		var user_id int
		var item string
		var modifiers json.RawMessage
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
	fmt.Println("Getting from redis")
	version, err := rdb.Get(ctx, fmt.Sprintf("restaurant:%d:active", restaurant_id)).Result()  //get the curr active version
	//if yes, return the cached json
	//if not, get menu id and version id then fetch items
	if err != nil{
		//fetch from db 
		row := conn.QueryRow(context.Background(), `SELECT v.id FROM version as v LEFT JOIN menu as m ON m.id=v.menu_id LEFT JOIN restaurant as r ON r.id=m.restaurant_id WHERE r.id=$1 AND v.is_active=true;`, restaurant_id)
		err := row.Scan(&version)
		fmt.Println(version)
		if err != nil{
			return []Item{}, fmt.Errorf("Error retreiving version %w", err)
		}
		//update redis with new one as well
		//row might need unmarshalling
		helper.UpdateRedis(ctx, rdb, fmt.Sprintf("restaurant:%d:active", restaurant_id), row)
	}
	//get menu
	json, err := rdb.Get(ctx, fmt.Sprintf("restaurant:%d:menu:v%s", restaurant_id, version)).Result()
	if err != nil{
		//fetch from postgres if doesn't exist
		
		rows, err := conn.Query(context.Background(), `SELECT id, version_id, name, description, price, category, modifiers FROM item WHERE version_id=$1;`, version)
		defer rows.Close()
		if err != nil{
			return []Item{}, fmt.Errorf("Error retreiving menu items %w", err)
		}
		res, err := pgx.CollectRows(rows, pgx.RowToStructByName[Item])
		if err != nil{
			return []Item{}, fmt.Errorf("Error converting menu items into struct %w", err)
		}
		//update redis with new menu items
		helper.UpdateRedis(ctx, rdb, fmt.Sprintf("restaurant:%d:menu:v%s", restaurant_id, version), res)
		return res, nil
	}
	fmt.Printf("%s", json)
	
	//turn json into Items

	return []Item{}, nil
}