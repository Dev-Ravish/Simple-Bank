package db

import (
	"context"
	"database/sql"
	"fmt"
)

type store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db:      db,
		Queries: New(db),
	}
}

func (store *store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			return fmt.Errorf("Transaction error %v, Rollback error %v", err, rErr)
		}
		return err
	}

	return tx.Commit()

}

var txKey = struct{}{}

type TransferTxParams struct {
	FromAccountId int64 `json:"transfer_from_account"`
	ToAccountId   int64 `json:"transfer_to_account"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer        Transfer    `json:"transfer"`
	FromAccount     Account     `json:"from_account"`
	ToAccount       Account     `json:"to_account"`
	FromTransaction Transaction `json:"from_transaction"`
	ToTransaction   Transaction `json:"to_transaction"`
}

func (store *store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	txName := ctx.Value(txKey)

	fmt.Println(txName)
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.TransferAmount(ctx, TransferAmountParams{
			TransferFromAccount: arg.FromAccountId,
			TransferToAccount:   arg.ToAccountId,
			Amount:              arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
		if err != nil {
			return err
		}
		fmt.Println(txName, " - before ac1, ", account1.Amount)

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:     arg.FromAccountId,
			Amount: account1.Amount - arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, " - after ac1, ", result.FromAccount.Amount)

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}
		fmt.Println(txName, " - before ac2, ", account2.Amount)

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:     arg.ToAccountId,
			Amount: account2.Amount + arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, " - after ac2, ", result.ToAccount.Amount)

		return nil
	})

	return result, err
}
