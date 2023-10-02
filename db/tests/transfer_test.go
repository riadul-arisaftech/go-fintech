package db_test

import (
	"context"
	"testing"

	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(user_id int64, t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		UserID:   int32(user_id),
		Balance:  200,
		Currency: "BDT",
	}

	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.UserID, arg.UserID)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomTransfer(t *testing.T, account1, account2 db.Account) db.TransferTxResponse {
	arg := db.CreateTransferParams{
		FromAccountID: int32(account1.ID),
		ToAccountID:   int32(account2.ID),
		Amount:        float64(utils.RandomMoney()),
	}

	tx, err := testQuery.TransferTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	// test transfer
	require.Equal(t, tx.Transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, tx.Transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, tx.Transfer.Amount, arg.Amount)
	// test entry
	require.Equal(t, tx.EntryIn.AccountID, arg.ToAccountID)
	require.Equal(t, tx.EntryIn.Amount, arg.Amount)

	require.Equal(t, tx.EntryOut.AccountID, arg.FromAccountID)
	require.Equal(t, tx.EntryOut.Amount, -1*arg.Amount)
	// test account
	require.Equal(t, tx.FromAccount.ID, account1.ID)
	require.Equal(t, tx.ToAccount.ID, account2.ID)

	require.Equal(t, tx.FromAccount.Balance, account1.Balance-arg.Amount)
	require.Equal(t, tx.ToAccount.Balance, account2.Balance+arg.Amount)
	return tx
}

func TestTransfer(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)

	account1 := createRandomAccount(user1.ID, t)
	account2 := createRandomAccount(user2.ID, t)

	arg := db.CreateTransferParams{
		FromAccountID: int32(account1.ID),
		ToAccountID:   int32(account2.ID),
		Amount:        10,
	}

	txResponseChan := make(chan db.TransferTxResponse)
	errorChan := make(chan error)
	count := 4
	for i := 0; i < count; i++ {
		go func() {
			tx, err := testQuery.TransferTx(context.Background(), arg)
			errorChan <- err
			txResponseChan <- tx
		}()
	}

	for i := 0; i < count; i++ {
		err := <-errorChan
		tx := <-txResponseChan

		require.NoError(t, err)
		require.NotEmpty(t, tx)
		// test transfer
		require.Equal(t, tx.Transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, tx.Transfer.ToAccountID, arg.ToAccountID)
		require.Equal(t, tx.Transfer.Amount, arg.Amount)
		// test entry
		require.Equal(t, tx.EntryIn.AccountID, arg.ToAccountID)
		require.Equal(t, tx.EntryIn.Amount, arg.Amount)

		require.Equal(t, tx.EntryOut.AccountID, arg.FromAccountID)
		require.Equal(t, tx.EntryOut.Amount, -1*arg.Amount)
		// test account
		// require.Equal(t, tx.FromAccount.ID, account1.ID)
		// require.Equal(t, tx.ToAccount.ID, account2.ID)

		// require.Equal(t, tx.FromAccount.Balance, account1.Balance-arg.Amount)
		// require.Equal(t, tx.ToAccount.Balance, account2.Balance+arg.Amount)
	}

	newAccount1, err := testQuery.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newAccount1)
	newAccount2, err := testQuery.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newAccount2)

	newAccount := float64(count * int(arg.Amount))
	require.Equal(t, newAccount1.Balance, account1.Balance-newAccount)
	require.Equal(t, newAccount2.Balance, account1.Balance+newAccount)
}
