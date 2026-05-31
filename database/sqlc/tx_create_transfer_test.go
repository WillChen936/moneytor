package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransferTx(t *testing.T) {
	ctx := context.Background()
	user := RandomUser(t)
	fromAccount := RandomAccount(t, user.ID)
	toAccount := RandomAccount(t, user.ID)
	amount := utils.RandomInt64Range(1, 100)

	transferCategory, err := testStore.CreateCategory(ctx, CreateCategoryParams{
		UserID:            user.ID,
		Name:              utils.RandomString(6),
		TransactionTypeID: 3, // TransactionTypeTransfer
	})
	require.NoError(t, err)

	arg := CreateTransferTxParams{
		Name:          utils.RandomString(6),
		Note:          utils.RandomString(6),
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		CategoryID:    transferCategory.ID,
		Amount:        amount,
	}

	result, err := testStore.CreateTransferTx(ctx, arg)
	require.NoError(t, err)

	entry := result.Entry
	require.Equal(t, entry.Name, arg.Name)
	require.Equal(t, entry.Note, arg.Note)
	require.Equal(t, entry.FromAccountID, arg.FromAccountID)
	require.True(t, entry.ToAccountID.Valid)
	require.Equal(t, entry.ToAccountID.Int64, arg.ToAccountID)
	require.Equal(t, entry.CategoryID, arg.CategoryID)
	require.Equal(t, entry.Amount, arg.Amount)

	updatedFrom, err := testStore.GetAccount(ctx, GetAccountParams{ID: fromAccount.ID, UserID: user.ID})
	require.NoError(t, err)
	require.Equal(t, fromAccount.Balance-amount, updatedFrom.Balance)

	updatedTo, err := testStore.GetAccount(ctx, GetAccountParams{ID: toAccount.ID, UserID: user.ID})
	require.NoError(t, err)
	require.Equal(t, toAccount.Balance+amount, updatedTo.Balance)
}
