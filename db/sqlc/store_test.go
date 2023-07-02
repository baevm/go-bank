package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	bob := createRandomAccount(t)  // Sender
	john := createRandomAccount(t) // Receiver

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Create n concurrent transfers
	// from Bob to John
	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: bob.ID,
				ToAccountId:   john.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}()
	}

	existed := make(map[int]bool)

	// check for results
	for i := 0; i < n; i++ {
		err := <-errs
		res := <-results

		require.NoError(t, err)
		require.NotEmpty(t, res)

		// Check transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.Amount, amount)
		require.Equal(t, bob.ID, transfer.FromAccountID)
		require.Equal(t, john.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		// Check from account entry (remove money)
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, bob.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// Check to account entry (add money)
		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, john.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, bob.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, john.ID)

		// check account balances
		diff1 := bob.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - john.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedBob, err := testQueries.GetAccount(context.Background(), bob.ID)
	require.NoError(t, err)

	updatedJohn, err := testQueries.GetAccount(context.Background(), john.ID)
	require.NoError(t, err)

	require.Equal(t, bob.Balance-int64(n)*amount, updatedBob.Balance)
	require.Equal(t, john.Balance+int64(n)*amount, updatedJohn.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	bob := createRandomAccount(t)  // Sender
	john := createRandomAccount(t) // Receiver

	n := 10
	amount := int64(10)
	errs := make(chan error)

	// Create n concurrent transfers
	// 5 transfers from Bob to John
	// 5 transfers from John to Bob
	for i := 0; i < n; i++ {
		fromAccountID := bob.ID
		toAccountId := john.ID

		if i%2 == 1 {
			fromAccountID = john.ID
			toAccountId = bob.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check for results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedBob, err := testQueries.GetAccount(context.Background(), bob.ID)
	require.NoError(t, err)

	updatedJohn, err := testQueries.GetAccount(context.Background(), john.ID)
	require.NoError(t, err)

	require.Equal(t, bob.Balance, updatedBob.Balance)
	require.Equal(t, john.Balance, updatedJohn.Balance)
}
