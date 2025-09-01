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
	category := RandomCategory(t)

	err := testQueries.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)
}

func TestGetCategory(t *testing.T) {
	category := RandomCategory(t)

	categoryGet, err := testQueries.GetCategory(context.Background(), category.ID)

	require.NoError(t, err)
	require.Equal(t, category.ID, categoryGet.ID)
	require.Equal(t, category.Name, categoryGet.Name)
	require.Equal(t, category.TransactionTypeID, categoryGet.TransactionTypeID)
	require.NotEmpty(t, category.CreatedAt, categoryGet.CreatedAt)
	require.WithinDuration(t, category.CreatedAt, categoryGet.CreatedAt, time.Second)
	require.NotEmpty(t, category.UpdatedAt, categoryGet.UpdatedAt)
	require.WithinDuration(t, category.UpdatedAt, categoryGet.UpdatedAt, time.Second)

	err = testQueries.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)
}

func TestUpdateCategory(t *testing.T) {
	t.Run("UpdateOnlyName", func(t *testing.T) {
		category := RandomCategory(t)

		newName := utils.RandomString(6)

		arg := UpdateCategoryParams{
			ID: category.ID,
			Name: pgtype.Text{
				String: newName,
				Valid:  true,
			},
			TransactionTypeID: pgtype.Int2{
				Valid: false,
			},
		}

		categoryUpdated, err := testQueries.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, newName, categoryUpdated.Name)
		require.Equal(t, category.TransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.After(category.UpdatedAt))

		err = testQueries.DeleteCategory(context.Background(), category.ID)
		require.NoError(t, err)
	})

	t.Run("UpdateOnlyTransactionID", func(t *testing.T) {
		category := RandomCategory(t)

		newTransactionTypeID := utils.RandomInt16Range(1, 3)

		arg := UpdateCategoryParams{
			ID: category.ID,
			Name: pgtype.Text{
				Valid: false,
			},
			TransactionTypeID: pgtype.Int2{
				Int16: newTransactionTypeID,
				Valid: true,
			},
		}

		categoryUpdated, err := testQueries.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, category.Name, categoryUpdated.Name)
		require.Equal(t, newTransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.After(category.UpdatedAt))

		err = testQueries.DeleteCategory(context.Background(), category.ID)
		require.NoError(t, err)
	})

	t.Run("UpdateAll", func(t *testing.T) {
		category := RandomCategory(t)

		newName := utils.RandomString(6)
		newTransactionTypeID := utils.RandomInt16Range(1, 3)

		arg := UpdateCategoryParams{
			ID: category.ID,
			Name: pgtype.Text{
				String: newName,
				Valid:  true,
			},
			TransactionTypeID: pgtype.Int2{
				Int16: newTransactionTypeID,
				Valid: true,
			},
		}

		categoryUpdated, err := testQueries.UpdateCategory(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, categoryUpdated)
		require.Equal(t, category.ID, categoryUpdated.ID)
		require.Equal(t, newName, categoryUpdated.Name)
		require.Equal(t, newTransactionTypeID, categoryUpdated.TransactionTypeID)
		require.WithinDuration(t, category.CreatedAt, categoryUpdated.CreatedAt, time.Second)
		require.True(t, categoryUpdated.UpdatedAt.After(category.UpdatedAt))

		err = testQueries.DeleteCategory(context.Background(), category.ID)
		require.NoError(t, err)
	})
}

func TestDeleteCategory(t *testing.T) {
	category := RandomCategory(t)

	errDelete := testQueries.DeleteCategory(context.Background(), category.ID)
	categoryGet, errGet := testQueries.GetCategory(context.Background(), category.ID)

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, categoryGet)
}

func RandomCategory(t *testing.T) Category {
	transactionType := RandomTransactionType(t)

	arg := CreateCategoryParams{
		Name:              utils.RandomString(6),
		TransactionTypeID: transactionType.ID,
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, category.ID)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.TransactionTypeID, category.TransactionTypeID)
	require.NotZero(t, category.CreatedAt)
	require.NotZero(t, category.UpdatedAt)

	return category
}
