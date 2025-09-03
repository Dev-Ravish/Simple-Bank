package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
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

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	txName := ctx.Value(txKey)

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

		if arg.FromAccountId > arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
			if err != nil {
				return err
			}
		}

		fmt.Println(txName, result.ToAccount.Balance, result.FromAccount.Balance)

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	fromAccountId int64,
	fromAmount int64,
	toAccountId int64,
	toAmount int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountId,
		Amount: fromAmount,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountId,
		Amount: toAmount,
	})
	return
}
