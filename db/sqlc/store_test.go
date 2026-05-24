package db

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := "10.00"

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for range make([]int, n) {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	for range make([]int, n) {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, "-"+amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		account1Balance, _ := decimal.NewFromString(account1.Balance)
		fromAccountBalance, _ := decimal.NewFromString(fromAccount.Balance)
		account2Balance, _ := decimal.NewFromString(toAccount.Balance)
		toAccountBalance, _ := decimal.NewFromString(toAccount.Balance)
		zero := decimal.NewFromInt(0)
		amountDecimal, _ := decimal.NewFromString(amount)

		diff1 := account1Balance.Sub(fromAccountBalance)
		diff2 := toAccountBalance.Sub(account2Balance)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1.GreaterThan(zero))
		require.True(t, diff1.Mod(amountDecimal).IsZero())

		k := diff1.Div(amountDecimal)
		require.True(t, k.GreaterThan(decimal.NewFromInt(0)) && k.LessThan(decimal.NewFromInt(int64(n+1))))
	}
}
