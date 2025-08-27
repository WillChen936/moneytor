package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var countTransactionType int32 = 3

func TestGetTransactionType(t *testing.T) {
	// Arrange
	id := int32(1)

	// Act
	transactionType, err := testQueries.GetTransactionType(context.Background(), id)

	// Assert
	require.NoError(t, err)
	require.Equal(t, id, transactionType.ID)
	require.NotEmpty(t, transactionType.Name)
}

func TestListTransactionType(t *testing.T) {
	// Arrange
	// Act
	transactionTypes, err := testQueries.ListTransactionTypes(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, countTransactionType, int32(len(transactionTypes)))
	for _, types := range transactionTypes {
		require.NotZero(t, types.ID)
		require.NotEmpty(t, types.Name)
	}
}

func RandomTransactionType() (TransactionType, error) {
	id := utils.RandomInt32Range(1, countTransactionType)
	return testQueries.GetTransactionType(context.Background(), id)
}
