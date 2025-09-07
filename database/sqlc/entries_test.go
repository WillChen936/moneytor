package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatEntry(t *testing.T) {
	testQueries := setupTestQueries(t)
	account := RandomAccount(t, testQueries)
	RandomEntry(t, testQueries, account.ID)
}

func TestListEntries(t *testing.T) {
	testQueries := setupTestQueries(t)
	account := RandomAccount(t, testQueries)
	var lastEntry Entry
	for i := 0; i < 10; i++ {
		lastEntry = RandomEntry(t, testQueries, account.ID)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    1,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, lastEntry.AccountID, account.ID)
	}
}

func RandomEntry(t *testing.T, testQueries *Queries, accountID int64) Entry {
	category := RandomCategory(t, testQueries)

	arg := CreateEntryParams{
		AccountID:  accountID,
		CategoryID: category.ID,
		Amount:     utils.RandomDecimalRange(100, 10000, 2),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry.ID)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.CategoryID, entry.CategoryID)
	require.True(t, arg.Amount.Equal(entry.Amount))
	require.NotEmpty(t, entry.CreatedAt)

	return entry
}
