package main

import (
	"fmt"
	"os"
	"log"

	"github.com/MatthewCCChang/Bento/backend/pkg/db/create"
	// "github.com/MatthewCCChang/Bento/backend/pkg/db/delete"
	"github.com/MatthewCCChang/Bento/backend/pkg/db/update"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")
	
	err := godotenv.Load()
	fmt.Printf("url is %s", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
	conn, err := create.CreateConnectionPool(10)
	if err != nil {
		fmt.Printf("Error creating connection pool: %v\n", err)
		return
	}
	fmt.Println("Successfully created connection pool")
	defer conn.Close()
	// rows, err := create.CreateDatabase(conn)
	// if err != nil {
	// 	fmt.Printf("Error creating database: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Created %d rows\n", rows)

	//deleting tables before creating them to avoid errors
	//err = delete.DeleteTables(conn)
	//creating tables
	// err = create.CreateTables(conn)
	//create indiv tables
	// _, err = create.CreateTable(conn, "users", "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, uuid TEXT NOT NULL UNIQUE, email TEXT, name TEXT, password TEXT")
	// if err != nil {
	// 	fmt.Printf("Error creating users table: %v\n", err)
	// 	return
	// }

	//inserting rows into tables
	// _, err = update.InsertIntoTable(conn, "users", []string{"uuid", "email", "name", "password"}, []interface{}{"fa7cf439-0b66-43cb-9a9e-6e7abe026b5a", "123test@gmail.com", "John Doe", "password123"})
	// if err != nil {
	// 	fmt.Printf("Error inserting into table: %v\n", err)
	// }
	// _, err = update.InsertIntoTable(conn, "restaurant", 
    // []string{"name", "address", "phone"}, 
    // []interface{}{"The Syntax Grill", "101 Binary Way", "555-0101"})
	// if err != nil {
	// 	fmt.Printf("Error inserting into restaurant: %v\n", err)
	// }
	// _, err = update.InsertIntoTable(conn, "menu", 
    // []string{"restaurant_id"}, 
    // []interface{}{1})
	// if err != nil {
	// 	fmt.Printf("Error inserting into menu: %v\n", err)
	// }

	_, err = update.InsertIntoTable(conn, "version", 
    []string{"menu_id", "s3_url", "is_active"}, 
    []interface{}{1, "https://s3.amazonaws.com/assets/menu-v1.json", true})
	if err != nil {
		fmt.Printf("Error inserting into version: %v\n", err)
	}
	_, err = update.InsertIntoTable(conn, "item", 
    []string{"version_id", "name", "description", "price", "category", "modifiers"}, 
    []interface{}{1, "Database Burger", "Highly indexed flavors", 12.99, "Entrees", `{"extra_cheese": true, "temp": "medium"}`})
	if err != nil {
		fmt.Printf("Error inserting into item: %v\n", err)
	}
}
