package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(100)

	resultC := make(chan TransferTxResult)
	errC := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			fmt.Println("starting go routine", i)
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			errC <- err
			resultC <- result
			fmt.Println("res", resultC)

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errC
		fmt.Println("recieved: ", err)

		require.NoError(t, err)

		result := <-resultC
		fmt.Println("res recieved: ", result)
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, transfer.TransferFromAccount, account1.ID)
		require.Equal(t, transfer.TransferToAccount, account2.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		toTransaction := result.ToTransaction
		require.NotEmpty(t, toTransaction)
		require.NotZero(t, toTransaction.ID)
		require.NotZero(t, toTransaction.CreatedAt)
		require.Equal(t, toTransaction.AccountID, account2.ID)
		require.Equal(t, toTransaction.Amount, amount)

		_, err = store.GetTransaction(context.Background(), toTransaction.ID)
		require.NoError(t, err)

		fromTransaction := result.FromTransaction
		require.NotEmpty(t, fromTransaction)
		require.NotZero(t, fromTransaction.ID)
		require.NotZero(t, fromTransaction.CreatedAt)
		require.Equal(t, fromTransaction.AccountID, account1.ID)
		require.Equal(t, fromTransaction.Amount, -amount)

		_, err = store.GetTransaction(context.Background(), fromTransaction.ID)
		require.NoError(t, err)

		// acc, errA := store.GetAccount(context.Background(), account1.ID)
		// require.NoError(t, errA)
		// require.Equal(t, account1.Amount-amount, acc.Amount)

		// acc, errA = store.GetAccount(context.Background(), account1.ID)
		// require.NoError(t, errA)
		// require.Equal(t, account2.Amount+amount, acc.Amount)

	}

}
