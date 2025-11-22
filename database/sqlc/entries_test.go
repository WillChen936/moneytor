package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatEntry(t *testing.T) {
	account := RandomAccount(t, testStore)
	RandomEntry(t, testStore, account.ID)
}

func TestListEntries(t *testing.T) {
	account := RandomAccount(t, testStore)

	for i := 0; i < 10; i++ {
		RandomEntry(t, testStore, account.ID)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 1,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListEntriesByAccountID(t *testing.T) {
	account1 := RandomAccount(t, testStore)
	account2 := RandomAccount(t, testStore)

	for i := 0; i < 10; i++ {
		RandomEntry(t, testStore, account1.ID)
		RandomEntry(t, testStore, account2.ID)
	}

	arg := ListEntriesByAccountIDParams{
		AccountID: account2.ID,
		Limit:     5,
		Offset:    1,
	}

	entries, err := testStore.ListEntriesByAccountID(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account2.ID)
	}
}

func RandomEntry(t *testing.T, testStore Store, accountID int64) Entry {
	category := RandomCategory(t, testStore)

	arg := CreateEntryParams{
		Name:       utils.RandomString(10),
		Note:       utils.RandomString(30),
		AccountID:  accountID,
		CategoryID: category.ID,
		Amount:     utils.RandomDecimalRange(100, 10000, 2),
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)

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
