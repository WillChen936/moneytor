package db

import (
	"context"
	"math/big"
	"moneytor/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	account := RandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	account := RandomAccount(t)

	accountGet, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account.ID, accountGet.ID)
	require.Equal(t, account.Owner, accountGet.Owner)
	require.Equal(t, account.CurrencyID, accountGet.CurrencyID)
	require.Equal(t, account.Balance, accountGet.Balance)
	require.WithinDuration(t, account.CreatedAt, accountGet.CreatedAt, time.Second)
	require.WithinDuration(t, account.UpdatedAt, accountGet.UpdatedAt, time.Second)

	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	t.Run("UpdateOnlyOwner", func(t *testing.T) {
		account := RandomAccount(t)

		newOwner := utils.RandomString(6)

		arg := UpdateAccountParams{
			ID: account.ID,
			Owner: pgtype.Text{
				String: newOwner,
				Valid:  true,
			},
			Balance: pgtype.Numeric{
				Valid: false,
			},
		}

		accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, accountUpdated)
		require.Equal(t, account.ID, accountUpdated.ID)
		require.Equal(t, newOwner, accountUpdated.Owner)
		require.Equal(t, account.Balance, accountUpdated.Balance)
		require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
		require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
		require.True(t, accountUpdated.UpdatedAt.After(account.UpdatedAt))

		err = testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)
	})

	t.Run("UpdateOnlyBalance", func(t *testing.T) {
		account := RandomAccount(t)

		arg := UpdateAccountParams{
			ID: account.ID,
			Owner: pgtype.Text{
				Valid: false,
			},
			Balance: pgtype.Numeric{
				Int:   big.NewInt(utils.RandomInt64Range(1000000, 100000000)),
				Exp:   -6,
				Valid: true,
			},
		}

		accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, accountUpdated)
		require.Equal(t, account.ID, accountUpdated.ID)
		require.Equal(t, account.Owner, accountUpdated.Owner)
		require.Equal(t, decimal.NewFromBigInt(arg.Balance.Int, arg.Balance.Exp), accountUpdated.Balance)
		require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
		require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
		require.True(t, accountUpdated.UpdatedAt.After(account.UpdatedAt))

		err = testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)
	})

	t.Run("UpdateAll", func(t *testing.T) {
		account := RandomAccount(t)

		newOwner := utils.RandomString(6)

		arg := UpdateAccountParams{
			ID: account.ID,
			Owner: pgtype.Text{
				String: newOwner,
				Valid:  true,
			},
			Balance: pgtype.Numeric{
				Int:   big.NewInt(utils.RandomInt64Range(1000000, 100000000)),
				Exp:   -6,
				Valid: true,
			},
		}

		accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, accountUpdated)
		require.Equal(t, account.ID, accountUpdated.ID)
		require.Equal(t, newOwner, accountUpdated.Owner)
		require.Equal(t, decimal.NewFromBigInt(arg.Balance.Int, arg.Balance.Exp), accountUpdated.Balance)
		require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
		require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
		require.True(t, accountUpdated.UpdatedAt.After(account.UpdatedAt))

		err = testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)
	})
}

func TestUpdateAccountBalance(t *testing.T) {
	account := RandomAccount(t)

	amount, err := decimal.NewFromString("12.34")
	require.NoError(t, err)

	arg := UpdateAccountBalanceParams{
		ID: account.ID,
		Amount: pgtype.Numeric{
			Int:   amount.Coefficient(),
			Exp:   amount.Exponent(),
			Valid: true,
		},
	}

	accountUpdated, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	expectedAmount := account.Balance.Add(amount)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)
	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Owner, accountUpdated.Owner)
	require.True(t, expectedAmount.Equal(accountUpdated.Balance))
	require.Equal(t, account.CurrencyID, accountUpdated.CurrencyID)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
	require.True(t, accountUpdated.UpdatedAt.After(account.UpdatedAt))

	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	account := RandomAccount(t)

	errDelete := testQueries.DeleteAccount(context.Background(), account.ID)
	accountGet, errGet := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, accountGet)
}

func RandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:      utils.RandomString(6),
		CurrencyID: RandomCurrency(t).ID,
		Balance: pgtype.Numeric{
			Int:   big.NewInt(utils.RandomInt64Range(1000000, 100000000)),
			Exp:   -6,
			Valid: true,
		},
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account.ID)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.CurrencyID, account.CurrencyID)
	require.Equal(t, decimal.NewFromBigInt(arg.Balance.Int, arg.Balance.Exp), account.Balance)
	require.NotEmpty(t, account.CreatedAt)
	require.NotEmpty(t, account.UpdatedAt)

	return account
}
