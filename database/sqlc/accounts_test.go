package db

import (
	"context"
	"moneytor/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	testQueries := setupTestQueries(t)
	RandomAccount(t, testQueries)
}

func TestGetAccount(t *testing.T) {
	testQueries := setupTestQueries(t)
	account := RandomAccount(t, testQueries)

	accountGet, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account.ID, accountGet.ID)
	require.Equal(t, account.Name, accountGet.Name)
	require.Equal(t, account.CurrencyID, accountGet.CurrencyID)
	require.Equal(t, account.Balance, accountGet.Balance)
	require.WithinDuration(t, account.CreatedAt, accountGet.CreatedAt, time.Second)
}

func TestUpdateAccountBalance(t *testing.T) {
	testQueries := setupTestQueries(t)
	account := RandomAccount(t, testQueries)

	amount, err := decimal.NewFromString("12.34")
	require.NoError(t, err)

	arg := UpdateAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	}

	accountUpdated, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	expectedAmount := account.Balance.Add(amount)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)
	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Name, accountUpdated.Name)
	require.True(t, expectedAmount.Equal(accountUpdated.Balance))
	require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
	require.True(t, accountUpdated.UpdatedAt.Valid)
	require.GreaterOrEqual(t, accountUpdated.UpdatedAt.Time, accountUpdated.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	testQueries := setupTestQueries(t)
	account := RandomAccount(t, testQueries)

	errDelete := testQueries.DeleteAccount(context.Background(), account.ID)
	accountGet, errGet := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, accountGet)
}

func RandomAccount(t *testing.T, testQueries *Queries) Account {
	arg := CreateAccountParams{
		Name:       utils.RandomString(6),
		CurrencyID: RandomCurrency(t, testQueries).ID,
		Balance:    utils.RandomDecimalRange(100, 10000, 2),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account.ID)
	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.CurrencyID, account.CurrencyID)
	require.True(t, arg.Balance.Equal(account.Balance))
	require.NotEmpty(t, account.CreatedAt)
	require.False(t, account.UpdatedAt.Valid)

	return account
}
