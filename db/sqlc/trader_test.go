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

	trader, err := testQueries.CreateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trader)

	require.Equal(t, arg.Account, trader.Account)
	require.Equal(t, arg.Balance, trader.Balance)
	require.Equal(t, arg.Currency, trader.Currency)

	require.NotZero(t, trader.ID)
	require.NotZero(t, trader.CreatedAt)

	return trader
}

func TestCreateTrader(t *testing.T) {
	createRandomTrader(t)
}

func TestGetTrader(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2, err := testQueries.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trader2)

	require.Equal(t, trader1.ID, trader2.ID)
	require.Equal(t, trader1.Account, trader2.Account)
	require.Equal(t, trader1.Balance, trader2.Balance)
	require.Equal(t, trader1.Currency, trader2.Currency)
	require.WithinDuration(t, trader1.CreatedAt, trader2.CreatedAt, time.Second)
}

func TestUpdateTrader(t *testing.T) {
	trader1 := createRandomTrader(t)

	arg := UpdateTraderParams{
		ID:      trader1.ID,
		Balance: util.RandomMoney(),
	}

	trader2, err := testQueries.UpdateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trader2)

	require.Equal(t, trader1.ID, trader2.ID)
	require.Equal(t, trader1.Account, trader2.Account)
	require.Equal(t, arg.Balance, trader2.Balance)
	require.Equal(t, trader1.Currency, trader2.Currency)
	require.WithinDuration(t, trader1.CreatedAt, trader2.CreatedAt, time.Second)
}

func TestDeleteTrader(t *testing.T) {
	trader1 := createRandomTrader(t)
	err := testQueries.DeleteTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	trader2, err := testQueries.GetTrader(context.Background(), trader1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trader2)
}

func TestListTraders(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTrader(t)
	}

	arg := ListTradersParams{
		Limit:  5,
		Offset: 0,
	}

	traders, err := testQueries.ListTraders(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, traders)
	require.Len(t, traders, 5)

	for _, trader := range traders {
		require.NotEmpty(t, trader)
	}
}
