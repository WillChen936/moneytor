package db

import (
	"context"
	"log"
	"moneytor/utils"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../config.json")
	if err != nil {
		log.Fatal(err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	defer connPool.Close()

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
