package delete

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DeleteTable(conn *pgxpool.Pool, name string) (int64, error) {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name)
	tag, err := conn.Exec(context.Background(), query)
	return tag.RowsAffected(), err
}

func DeleteTables(conn *pgxpool.Pool) error {
	tables := []string{"menu", "version", "item", "users", "restaurant"}
	for _, name := range tables {
		fmt.Printf("Deleting table %s\n", name)
		rows, err := DeleteTable(conn, name)	
		if err != nil {
			log.Printf("Error deleting table %s: %v\n", name, err)
			fmt.Println("Continuing to next table...")
			continue
		}
		fmt.Printf("Deleted %d rows\n", rows)
	}
	return nil
}