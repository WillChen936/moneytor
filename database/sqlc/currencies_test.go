package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var countCurrencies int16 = 5

func TestGetCurrency(t *testing.T) {
	testQueries := setupTestQueries(t)
	RandomCurrency(t, testQueries)
}

func TestListCurrencies(t *testing.T) {
	testQueries := setupTestQueries(t)
	currencies, err := testQueries.ListCurrencies(context.Background())

	require.NoError(t, err)
	require.Equal(t, countCurrencies, int16(len(currencies)))
	for _, currency := range currencies {
		require.NotZero(t, currency.ID)
		require.NotEmpty(t, currency.Code)
	}
}

func RandomCurrency(t *testing.T, q *Queries) Currency {
	id := utils.RandomInt16Range(1, countCurrencies)

	currency, err := q.GetCurrency(context.Background(), id)

	require.NoError(t, err)
	require.NotEmpty(t, currency)
	require.Equal(t, id, currency.ID)
	require.NotEmpty(t, currency.Code)

	return currency
}
