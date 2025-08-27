package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCurrency(t *testing.T) {
	expected := Currency{
		ID:           1,
		CurrencyCode: "TWD",
	}

	currency, err := testQueries.GetCurrency(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, currency)
	require.Equal(t, expected, currency)
}
