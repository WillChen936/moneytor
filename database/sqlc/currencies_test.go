package db

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestGetCurrency(t *testing.T) {
	expected := Currency{
		ID:           1,
		CurrencyCode: "TWD",
	}

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, "postgres://root:pass.123@localhost:5432/moneytor")
	if err != nil {
		log.Fatal(err)
	}

	queries := New(connPool)

	currency, err := queries.GetCurrency(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, currency)
	require.Equal(t, expected, currency)
}
