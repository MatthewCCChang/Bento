package update

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MatthewCCChang/Bento/backend/pkg/db/helper"
)

//UpdateUser()
func UpdateUser(conn *pgxpool.Pool, columns []string, values []interface{}, id int) (map[string]interface{}, error) {
	res, err := helper.UpdateTable(conn, "users", columns, []string{"id"}, values, []any{id}, []string{"name", "email"}, true)
	if err != nil{
		return map[string]interface{}{}, err
	}
	return res, nil
}

//updateItem()
func UpdateItem(conn *pgxpool.Pool, columns []string, values []interface{}, id int) (map[string]interface{}, error) {
	res, err := helper.UpdateTable(conn, "item", columns, []string{"id"}, values, []any{id}, []string{"name", "version", "price", "category", "modifiers"}, true)
	if err != nil{
		return map[string]interface{}{}, err
	}
	fmt.Println(res)
	return res, nil
}

//updateOrder()
func UpdateOrderItem(conn *pgxpool.Pool, columns []string, values []interface{}, id int) (map[string]interface{}, error) {
	res, err := helper.UpdateTable(conn, "order_items", columns,[]string{"id"}, values, []any{id}, []string{"item", "user_id", "price", "modifiers"}, true)
	if err != nil{
		return map[string]interface{}{}, err
	}
	fmt.Println(res)
	return res, nil
}

//updateMenu()  updated_at
func UpdateMenu(conn *pgxpool.Pool, columns []string, values []interface{}, id int) (map[string]interface{}, error) {
	res, err := helper.UpdateTable(conn, "menu", columns,[]string{"id"}, values, []any{id}, []string{"updated_at", "restaurant_id"}, true)
	if err != nil{
		return map[string]interface{}{}, err 
	}
	fmt.Println(res)
	return res, nil
}

//updateRestaurant()
func UpdateRestaurant(conn *pgxpool.Pool, columns []string, values []interface{}, id int) (map[string]interface{}, error) {
	res, err := helper.UpdateTable(conn, "item", columns,[]string{"id"}, values, []any{id}, []string{"name", "address", "telephone"}, true)
	if err != nil{
		return map[string]interface{}{}, err 
	}
	fmt.Println(res)
	return res, nil
}
