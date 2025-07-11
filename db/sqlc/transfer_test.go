package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/chicopalermo/simple-api-golang/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	deletedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedTransfer)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: util.RandomMoney(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)

	require.Equal(t, updatedTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, updatedTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, updatedTransfer.Amount)
	require.WithinDuration(t, updatedTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}
