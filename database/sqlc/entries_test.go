package db

import (
	"testing"
)

func TestCreatEntry(t *testing.T) {
	RandomEntry(t)
}

func RandomEntry(t *testing.T) {
	// category := RandomCategory(t)

	// arg := CreateEntryParams{
	// 	AccountID:  utils.RandomInt64Range(1, math.MaxInt64),
	// 	CategoryID: category.ID,
	// 	Amount: pgtype.Numeric{
	// 		Int:   big.NewInt(utils.RandomInt64Range(100000, math.MaxInt64-1)),
	// 		Exp:   -6,
	// 		Valid: true,
	// 	},
	// }

	// entry, err := testQueries.CreateEntry(context.Background(), arg)

	// require.NoError(t, err)
	// require.NotEmpty(t, entry.ID)
	// require.Equal(t, arg.AccountID, entry.AccountID)
	// require.Equal(t, arg.CategoryID, entry.CategoryID)
	// require.Equal(t, arg.Amount, entry.Amount)
	// require.NotEmpty(t, entry.CreatedAt)
}
