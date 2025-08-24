package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), "postgres://root:pass.123@localhost:5432/moneytor")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer connPool.Close()

	var id int
	var code string
	err = connPool.QueryRow(context.Background(), "SELECT * FROM currencies WHERE id=$1", 1).Scan(&id, &code)
	if err != nil {
		fmt.Printf("Query ro failed: %v", err)
		os.Exit(1)
	}

	fmt.Println(id, code)
}
