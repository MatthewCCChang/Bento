package main

import (
	"fmt"
	"os"
	"log"

	"github.com/MatthewCCChang/Bento/backend/internal/db/create"
	// "github.com/MatthewCCChang/Bento/backend/internal/db/delete"
	"github.com/MatthewCCChang/Bento/backend/internal/db/update"
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

	//inserting rows into tables
	_, err = update.InsertIntoTable(conn, "user", []string{"uuid", "email", "name", "password"}, []interface{}{"fa7cf439-0b66-43cb-9a9e-6e7abe026b5a", "123test@gmail.com", "John Doe", "password123"})
}
