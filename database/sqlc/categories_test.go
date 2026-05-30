package db

import (
	"context"
	"moneytor/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	user := RandomUser(t)
	RandomCategory(t, user.ID)
}

func TestGetCategory(t *testing.T) {
	user := RandomUser(t)
	category := RandomCategory(t, user.ID)

	categoryGet, err := testStore.GetCategory(context.Background(), GetCategoryParams{
		ID:     category.ID,
		UserID: category.UserID,
	})

	require.NoError(t, err)
	require.Equal(t, category.ID, categoryGet.ID)
	require.Equal(t, category.Name, categoryGet.Name)
	require.Equal(t, category.TransactionTypeID, categoryGet.TransactionTypeID)
	require.WithinDuration(t, category.CreatedAt, categoryGet.CreatedAt, time.Second)
	require.False(t, category.UpdatedAt.Valid)
}

func TestListCategories(t *testing.T) {
	user := RandomUser(t)
	for i := 0; i < 10; i++ {
		RandomCategory(t, user.ID)
	}

	arg := ListCategoriesParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}

	categories, err := testStore.ListCategories(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.Len(t, categories, int(arg.Limit))
}

func TestUpdateCategory(t *testing.T) {
	t.Run("UpdateOnlyName", func(t *testing.T) {
		user := RandomUser(t)
		category := RandomCategory(t, user.ID)
		newName := utils.RandomString(6)

		arg := UpdateCategoryParams{
			ID:     category.ID,
			UserID: category.UserID,
			Name: pgtype.Text{
				String: newName,
				Valid:  true,
			},
			TransactionTypeID: pgtype.Int2{Valid: false},
		}

		categoryUpdated, err := testStore.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, newName, categoryUpdated.Name)
		require.Equal(t, category.TransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})

	t.Run("UpdateOnlyTransactionID", func(t *testing.T) {
		user := RandomUser(t)
		category := RandomCategory(t, user.ID)
		newTransactionTypeID := utils.RandomInt16Range(1, 3)

		arg := UpdateCategoryParams{
			ID:     category.ID,
			UserID: category.UserID,
			Name:   pgtype.Text{Valid: false},
			TransactionTypeID: pgtype.Int2{
				Int16: newTransactionTypeID,
				Valid: true,
			},
		}

		categoryUpdated, err := testStore.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, category.Name, categoryUpdated.Name)
		require.Equal(t, newTransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})

	t.Run("UpdateAll", func(t *testing.T) {
		user := RandomUser(t)
		category := RandomCategory(t, user.ID)
		newName := utils.RandomString(6)
		newTransactionTypeID := utils.RandomInt16Range(1, 3)

		arg := UpdateCategoryParams{
			ID:     category.ID,
			UserID: category.UserID,
			Name: pgtype.Text{
				String: newName,
				Valid:  true,
			},
			TransactionTypeID: pgtype.Int2{
				Int16: newTransactionTypeID,
				Valid: true,
			},
		}

		categoryUpdated, err := testStore.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, newName, categoryUpdated.Name)
		require.Equal(t, newTransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})
}

func TestDeleteCategory(t *testing.T) {
	user := RandomUser(t)
	category := RandomCategory(t, user.ID)

	errDelete := testStore.DeleteCategory(context.Background(), DeleteCategoryParams{
		ID:     category.ID,
		UserID: category.UserID,
	})
	categoryGet, errGet := testStore.GetCategory(context.Background(), GetCategoryParams{
		ID:     category.ID,
		UserID: category.UserID,
	})

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, categoryGet)
}

func RandomCategory(t *testing.T, userID int64) Category {
	transactionType := RandomTransactionType(t)

	arg := CreateCategoryParams{
		UserID:            userID,
		Name:              utils.RandomString(6),
		TransactionTypeID: transactionType.ID,
	}

	category, err := testStore.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, category.ID)
	require.Equal(t, arg.UserID, category.UserID)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.TransactionTypeID, category.TransactionTypeID)
	require.NotZero(t, category.CreatedAt)
	require.False(t, category.UpdatedAt.Valid)

	return category
}
