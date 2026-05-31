package db

import (
	"context"
	"moneytor/utils"
	"sync"
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

// TestCreateTransferTxConcurrent verifies that concurrent transfers between the same two accounts
// do not cause deadlocks and produce correct final balances.
func TestCreateTransferTxConcurrent(t *testing.T) {
	ctx := context.Background()
	user := RandomUser(t)
	accountA := RandomAccount(t, user.ID)
	accountB := RandomAccount(t, user.ID)

	transferCategory, err := testStore.CreateCategory(ctx, CreateCategoryParams{
		UserID:            user.ID,
		Name:              utils.RandomString(6),
		TransactionTypeID: 3, // TransactionTypeTransfer
	})
	require.NoError(t, err)

	const n = 5
	amount := int64(10)

	errs := make(chan error, n*2)
	var wg sync.WaitGroup

	// n goroutines transfer A→B, n goroutines transfer B→A concurrently
	for i := 0; i < n; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, err := testStore.CreateTransferTx(ctx, CreateTransferTxParams{
				Name:          utils.RandomString(6),
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				CategoryID:    transferCategory.ID,
				Amount:        amount,
			})
			errs <- err
		}()
		go func() {
			defer wg.Done()
			_, err := testStore.CreateTransferTx(ctx, CreateTransferTxParams{
				Name:          utils.RandomString(6),
				FromAccountID: accountB.ID,
				ToAccountID:   accountA.ID,
				CategoryID:    transferCategory.ID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	// Net transfers cancel out — balances should be unchanged
	finalA, err := testStore.GetAccount(ctx, GetAccountParams{ID: accountA.ID, UserID: user.ID})
	require.NoError(t, err)
	require.Equal(t, accountA.Balance, finalA.Balance)

	finalB, err := testStore.GetAccount(ctx, GetAccountParams{ID: accountB.ID, UserID: user.ID})
	require.NoError(t, err)
	require.Equal(t, accountB.Balance, finalB.Balance)
}
