package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var connPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	connPool, err = pgxpool.New(context.Background(), "postgres://root:pass.123@localhost:5432/moneytor")
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
