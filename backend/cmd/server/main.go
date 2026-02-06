package main

import (
	"fmt"
	"os"
	"log"

	"github.com/MatthewCCChang/Bento/backend/internal/db/create"
	"github.com/MatthewCCChang/Bento/backend/internal/db/delete"
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
	err = delete.DeleteTables(conn)
	err = create.CreateTables(conn)
}
