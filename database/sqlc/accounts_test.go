package db

import (
	"context"
	"moneytor/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	RandomAccount(t, testStore)
}

func TestGetAccount(t *testing.T) {
	account := RandomAccount(t, testStore)

	accountGet, err := testStore.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account.ID, accountGet.ID)
	require.Equal(t, account.Name, accountGet.Name)
	require.Equal(t, account.CurrencyID, accountGet.CurrencyID)
	require.Equal(t, account.Balance, accountGet.Balance)
	require.WithinDuration(t, account.CreatedAt, accountGet.CreatedAt, time.Second)
}

func TestUpdateAccountBalance(t *testing.T) {
	account := RandomAccount(t, testStore)

	amount := utils.RandomInt64Range(-100, 100)

	arg := UpdateAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	}

	accountUpdated, err := testStore.UpdateAccountBalance(context.Background(), arg)
	expectedAmount := account.Balance + amount

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)
	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Name, accountUpdated.Name)
	require.Equal(t, expectedAmount, accountUpdated.Balance)
	require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
	require.True(t, accountUpdated.UpdatedAt.Valid)
	require.GreaterOrEqual(t, accountUpdated.UpdatedAt.Time, accountUpdated.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		RandomAccount(t, testStore)
	}

	limit := int32(5)
	offset := int32(0)

	arg := ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, int(limit))
}

func TestDeleteAccount(t *testing.T) {
	account := RandomAccount(t, testStore)

	errDelete := testStore.DeleteAccount(context.Background(), account.ID)
	accountGet, errGet := testStore.GetAccount(context.Background(), account.ID)

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, accountGet)
}

func RandomAccount(t *testing.T, testStore Store) Account {
	arg := CreateAccountParams{
		Name:       utils.RandomString(6),
		CurrencyID: RandomCurrency(t, testStore).ID,
		Balance:    utils.RandomInt64Range(100, 10000),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account.ID)
	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.CurrencyID, account.CurrencyID)
	require.Equal(t, arg.Balance, account.Balance)
	require.NotEmpty(t, account.CreatedAt)
	require.False(t, account.UpdatedAt.Valid)

	return account
}
