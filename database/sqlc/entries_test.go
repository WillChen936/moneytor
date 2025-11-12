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

	for i := 0; i < 10; i++ {
		RandomEntry(t, testQueries, account.ID)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 1,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListEntriesByAccountID(t *testing.T) {
	testQueries := setupTestQueries(t)
	account1 := RandomAccount(t, testQueries)
	account2 := RandomAccount(t, testQueries)

	for i := 0; i < 10; i++ {
		RandomEntry(t, testQueries, account1.ID)
		RandomEntry(t, testQueries, account2.ID)
	}

	arg := ListEntriesByAccountIDParams{
		AccountID: account2.ID,
		Limit:     5,
		Offset:    1,
	}

	entries, err := testQueries.ListEntriesByAccountID(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account2.ID)
	}
}

func RandomEntry(t *testing.T, testQueries *Queries, accountID int64) Entry {
	category := RandomCategory(t, testQueries)

	arg := CreateEntryParams{
		Name:       utils.RandomString(10),
		Note:       utils.RandomString(30),
		AccountID:  accountID,
		CategoryID: category.ID,
		Amount:     utils.RandomDecimalRange(100, 10000, 2),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry.ID)
	require.Equal(t, arg.Name, entry.Name)
	require.Equal(t, arg.Note, entry.Note)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.CategoryID, entry.CategoryID)
	require.True(t, arg.Amount.Equal(entry.Amount))
	require.NotEmpty(t, entry.CreatedAt)

	return entry
}
