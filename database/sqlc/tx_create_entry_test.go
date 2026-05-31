package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntryTx(t *testing.T) {
	ctx := context.Background()
	user := RandomUser(t)
	account := RandomAccount(t, user.ID)
	amount := utils.RandomInt64Range(-100, 100)

	arg := CreateEntryTxParams{
		Name:          utils.RandomString(6),
		Note:          utils.RandomString(6),
		FromAccountID: account.ID,
		CategoryID:    RandomCategory(t, user.ID).ID,
		Amount:        amount,
	}

	result, err := testStore.CreateEntryTx(ctx, arg)
	require.NoError(t, err)

	entryCreated := result.Entry
	amountUpdated, err := testStore.GetAccount(ctx, GetAccountParams{ID: account.ID, UserID: account.UserID})
	require.NoError(t, err)

	require.Equal(t, entryCreated.Name, arg.Name)
	require.Equal(t, entryCreated.Note, arg.Note)
	require.Equal(t, entryCreated.FromAccountID, arg.FromAccountID)
	require.Equal(t, entryCreated.CategoryID, arg.CategoryID)
	require.Equal(t, entryCreated.Amount, arg.Amount)
	require.Equal(t, amountUpdated.Balance, account.Balance+amount)
}
