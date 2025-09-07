package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var countTransactionType int16 = 3

func TestGetTransactionType(t *testing.T) {
	testQueries := setupTestQueries(t)
	RandomTransactionType(t, testQueries)
}

func TestListTransactionTypes(t *testing.T) {
	testQueries := setupTestQueries(t)
	transactionTypes, err := testQueries.ListTransactionTypes(context.Background())

	require.NoError(t, err)
	require.Equal(t, countTransactionType, int16(len(transactionTypes)))
	for _, types := range transactionTypes {
		require.NotZero(t, types.ID)
		require.NotEmpty(t, types.Name)
	}
}

func RandomTransactionType(t *testing.T, testQueries *Queries) TransactionType {
	id := utils.RandomInt16Range(1, countTransactionType)

	transactionType, err := testQueries.GetTransactionType(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, id, transactionType.ID)
	require.NotEmpty(t, transactionType.Name)

	return transactionType
}
