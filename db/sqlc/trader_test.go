package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/YuanData/aquatrade/util"
	"github.com/stretchr/testify/require"
)

func createRandomTrader(t *testing.T) Trader {
	arg := CreateTraderParams{
		Account:  util.RandomAccount(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Account, account.Account)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateTrader(t *testing.T) {
	createRandomTrader(t)
}

func TestGetTrader(t *testing.T) {
	account1 := createRandomTrader(t)
	account2, err := testQueries.GetTrader(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Account, account2.Account)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateTrader(t *testing.T) {
	account1 := createRandomTrader(t)

	arg := UpdateTraderParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Account, account2.Account)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteTrader(t *testing.T) {
	account1 := createRandomTrader(t)
	err := testQueries.DeleteTrader(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetTrader(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListTraders(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTrader(t)
	}

	arg := ListTradersParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListTraders(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
