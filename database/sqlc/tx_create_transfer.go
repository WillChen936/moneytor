package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTransferTxParams struct {
	Name          string
	Note          string
	FromAccountID int64
	ToAccountID   int64
	CategoryID    int64
	Amount        int64
}

type CreateTransferTxResult struct {
	Entry       Entry
	FromAccount Account
	ToAccount   Account
}

func (store *SQLStore) CreateTransferTx(ctx context.Context, arg CreateTransferTxParams) (CreateTransferTxResult, error) {
	var result CreateTransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			Name:          arg.Name,
			Note:          arg.Note,
			FromAccountID: arg.FromAccountID,
			ToAccountID:   pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
			CategoryID:    arg.CategoryID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update both accounts in ID order to prevent deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{ID: arg.FromAccountID, Amount: -arg.Amount})
			if err != nil {
				return err
			}
			result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
		} else {
			result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{ID: arg.FromAccountID, Amount: -arg.Amount})
		}
		return err
	})

	return result, err
}
