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
	testQueries := setupTestQueries(t)
	category := RandomCategory(t, testQueries)

	err := testQueries.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)
}

func TestGetCategory(t *testing.T) {
	testQueries := setupTestQueries(t)
	category := RandomCategory(t, testQueries)

	categoryGet, err := testQueries.GetCategory(context.Background(), category.ID)

	require.NoError(t, err)
	require.Equal(t, category.ID, categoryGet.ID)
	require.Equal(t, category.Name, categoryGet.Name)
	require.Equal(t, category.TransactionTypeID, categoryGet.TransactionTypeID)
	require.NotEmpty(t, category.CreatedAt, categoryGet.CreatedAt)
	require.WithinDuration(t, category.CreatedAt, categoryGet.CreatedAt, time.Second)
	require.False(t, category.UpdatedAt.Valid)
}

func TestUpdateCategory(t *testing.T) {
	t.Run("UpdateOnlyName", func(t *testing.T) {
		testQueries := setupTestQueries(t)
		category := RandomCategory(t, testQueries)

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
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})

	t.Run("UpdateOnlyTransactionID", func(t *testing.T) {
		testQueries := setupTestQueries(t)
		category := RandomCategory(t, testQueries)

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
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})

	t.Run("UpdateAll", func(t *testing.T) {
		testQueries := setupTestQueries(t)
		category := RandomCategory(t, testQueries)

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
		require.True(t, categoryUpdated.UpdatedAt.Valid)
		require.GreaterOrEqual(t, categoryUpdated.UpdatedAt.Time, categoryUpdated.CreatedAt)
	})
}

func TestDeleteCategory(t *testing.T) {
	testQueries := setupTestQueries(t)
	category := RandomCategory(t, testQueries)

	errDelete := testQueries.DeleteCategory(context.Background(), category.ID)
	categoryGet, errGet := testQueries.GetCategory(context.Background(), category.ID)

	require.NoError(t, errDelete)
	require.Error(t, errGet)
	require.EqualError(t, errGet, pgx.ErrNoRows.Error())
	require.Empty(t, categoryGet)
}

func RandomCategory(t *testing.T, testQueries *Queries) Category {
	transactionType := RandomTransactionType(t, testQueries)

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
	require.False(t, category.UpdatedAt.Valid)

	return category
}
