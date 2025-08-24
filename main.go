package main

import (
	"context"
	"fmt"
	db "moneytor/database/sqlc"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, "postgres://root:pass.123@localhost:5432/moneytor")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer connPool.Close()

	queries := db.New(connPool)
	currency, err := queries.GetCurrency(ctx, 1)
	if err != nil {
		fmt.Printf("Unable to get currency_id = 1")
		os.Exit(1)
	}

	fmt.Println(currency)
}
