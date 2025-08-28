package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Amount:   util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.Nil(t, err)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Amount, account.Amount)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
