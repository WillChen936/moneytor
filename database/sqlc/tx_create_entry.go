package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEntryTxParams struct {
	Name       string
	Note       string
	AccountID  int64
	CategoryID int64
	Amount     int64
}

type CreateEntryTxResult struct {
	Entry   Entry
	Account Account
}

// CreateEntryTx creates an entry and updates the account balance (balance + amount).
// Negative balance is allowed (e.g. credit card accounts may be overdrawn), so no CHECK (balance >= 0) is applied.
func (store *SQLStore) CreateEntryTx(ctx context.Context, arg CreateEntryTxParams) (CreateEntryTxResult, error) {
	var result CreateEntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			Name:          arg.Name,
			Note:          arg.Note,
			FromAccountID: arg.AccountID,
			ToAccountID:   pgtype.Int8{Valid: false},
			CategoryID:    arg.CategoryID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.Account, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
