package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntryTx(t *testing.T) {
	ctx := context.Background()
	account := RandomAccount(t, testStore)
	amount := utils.RandomDecimalRange(100, 10000, 2).Neg()

	arg := CreateEntryTxParams{
		Name:       utils.RandomString(6),
		Note:       utils.RandomString(6),
		AccountID:  account.ID,
		CategoryID: RandomCategory(t, testStore).ID,
		Amount:     amount,
	}

	result, err := testStore.CreateEntryTx(ctx, arg)
	require.NoError(t, err)

	entryCreated := result.Entry
	amountUpdated, err := testStore.GetAccount(ctx, account.ID)
	require.NoError(t, err)

	require.Equal(t, entryCreated.Name, arg.Name)
	require.Equal(t, entryCreated.Note, arg.Note)
	require.Equal(t, entryCreated.AccountID, arg.AccountID)
	require.Equal(t, entryCreated.CategoryID, arg.CategoryID)
	require.True(t, entryCreated.Amount.Equal(arg.Amount))
	require.True(t, amountUpdated.Balance.Equal(account.Balance.Add(amount)))
}
