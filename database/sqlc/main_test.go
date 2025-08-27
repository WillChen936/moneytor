package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), "postgres://root:pass.123@localhost:5432/moneytor")
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
}
