package db

import (
	"context"
	"log"
	"moneytor/utils"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var connPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../config.json")
	if err != nil {
		log.Fatal(err)
	}

	connPool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	defer connPool.Close()

	os.Exit(m.Run())
}

func setupTestQueries(t *testing.T) *Queries {
	ctx := context.Background()

	tx, err := connPool.Begin(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		tx.Rollback(ctx)
	})

	return New(tx)
}
