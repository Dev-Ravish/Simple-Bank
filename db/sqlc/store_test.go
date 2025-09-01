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
	amount := int64(5)

	resultC := make(chan TransferTxResult)
	errC := make(chan error)

	for i := 0; i < n; i++ {
		txValue := fmt.Sprintf("transacction - %d ", i+1)
		go func() {
			cntx := context.WithValue(context.Background(), txKey, txValue)
			result, err := store.TransferTx(cntx, TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			resultC <- result
			errC <- err

		}()
	}

	existed := map[int]bool{}
	for i := 0; i < n; i++ {
		result := <-resultC
		err := <-errC

		require.NotEmpty(t, result)
		require.NoError(t, err)

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

		// fmt.Println(i, "acc1 - ", acc1.Balance, "********* acc2 - ", acc2.Balance)
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err1 := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err1)

	updatedAccount2, err2 := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err2)

	diff1 := account1.Balance - updatedAccount1.Balance
	diff2 := updatedAccount2.Balance - account2.Balance
	require.Equal(t, diff1, diff2)
	require.Equal(t, diff1, 5*amount)
}
