package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var countCurrencies int16 = 5

func TestGetCurrency(t *testing.T) {
	expected := Currency{
		ID:   1,
		Code: "TWD",
	}

	currency, err := testQueries.GetCurrency(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, currency)
	require.Equal(t, expected, currency)
}

func TestListCurrencies(t *testing.T) {
	// Arrange
	// Act
	currencies, err := testQueries.ListCurrencies(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, countCurrencies, int16(len(currencies)))
	for _, currency := range currencies {
		require.NotZero(t, currency.ID)
		require.NotEmpty(t, currency.Code)
	}
}

func RandomCurrency() (Currency, error) {
	id := utils.RandomInt16Range(1, countCurrencies)
	return testQueries.GetCurrency(context.Background(), id)
}
