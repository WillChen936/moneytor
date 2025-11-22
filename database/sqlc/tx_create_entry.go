package db

import (
	"context"

	"github.com/shopspring/decimal"
)

type CreateEntryTxParams struct {
	Name       string
	Note       string
	AccountID  int64
	CategoryID int64
	Amount     decimal.Decimal
}

type CreateEntryTxResult struct {
	Entry   Entry
	Account Account
}

func (store *SQLStore) CreateEntryTx(ctx context.Context, arg CreateEntryTxParams) (CreateEntryTxResult, error) {
	var result CreateEntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams(arg))
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
