package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var countTransactionType int16 = 3

func TestGetTransactionType(t *testing.T) {
	// Arrange
	id := int16(1)

	// Act
	transactionType, err := testQueries.GetTransactionType(context.Background(), id)

	// Assert
	require.NoError(t, err)
	require.Equal(t, id, transactionType.ID)
	require.NotEmpty(t, transactionType.Name)
}

func TestListTransactionTypes(t *testing.T) {
	// Arrange
	// Act
	transactionTypes, err := testQueries.ListTransactionTypes(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, countTransactionType, int16(len(transactionTypes)))
	for _, types := range transactionTypes {
		require.NotZero(t, types.ID)
		require.NotEmpty(t, types.Name)
	}
}

func RandomTransactionType() (TransactionType, error) {
	id := utils.RandomInt16Range(1, countTransactionType)
	return testQueries.GetTransactionType(context.Background(), id)
}
