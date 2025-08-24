package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://root:pass.123@localhost:5432/moneytor")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var id int
	var code string
	err = conn.QueryRow(context.Background(), "SELECT * FROM currencies WHERE id=$1", 10).Scan(&id, &code)
	if err != nil {
		fmt.Printf("Query ro failed: %v", err)
		os.Exit(1)
	}

	fmt.Println(id, code)
}
