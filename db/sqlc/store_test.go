package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentTx(t *testing.T) {
	store := NewStore(testDB)

	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan PaymentTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PaymentTx(context.Background(), PaymentTxParams{
				FromTraderID: trader1.ID,
				ToTraderID:   trader2.ID,
				Amount:       amount,
			})

			errs <- err
			results <- result
		}()
	}

	// verify results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// verify payment
		payment := result.Payment
		require.NotEmpty(t, payment)
		require.Equal(t, trader1.ID, payment.FromTraderID)
		require.Equal(t, trader2.ID, payment.ToTraderID)
		require.Equal(t, amount, payment.Amount)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)

		_, err = store.GetPayment(context.Background(), payment.ID)
		require.NoError(t, err)

		// verify records
		fromRecord := result.FromRecord
		require.NotEmpty(t, fromRecord)
		require.Equal(t, trader1.ID, fromRecord.TraderID)
		require.Equal(t, -amount, fromRecord.Amount)
		require.NotZero(t, fromRecord.ID)
		require.NotZero(t, fromRecord.CreatedAt)

		_, err = store.GetRecord(context.Background(), fromRecord.ID)
		require.NoError(t, err)

		toRecord := result.ToRecord
		require.NotEmpty(t, toRecord)
		require.Equal(t, trader2.ID, toRecord.TraderID)
		require.Equal(t, amount, toRecord.Amount)
		require.NotZero(t, toRecord.ID)
		require.NotZero(t, toRecord.CreatedAt)

		_, err = store.GetRecord(context.Background(), toRecord.ID)
		require.NoError(t, err)

		// verify traders
		fromTrader := result.FromTrader
		// require.NotEmpty(t, fromTrader)
		require.Equal(t, trader1.ID, fromTrader.ID)

		toTrader := result.ToTrader
		// require.NotEmpty(t, toTrader)
		require.Equal(t, trader2.ID, toTrader.ID)

		diff1 := trader1.Balance - fromTrader.Balance
		diff2 := toTrader.Balance - trader2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// verify the final updated result
	updatedTrader1, err := store.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	updatedTrader2, err := store.GetTrader(context.Background(), trader2.ID)
	require.NoError(t, err)

	require.Equal(t, trader1.Balance-int64(n)*amount, updatedTrader1.Balance)
	require.Equal(t, trader2.Balance+int64(n)*amount, updatedTrader2.Balance)
}

func TestPaymentTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromTraderID := trader1.ID
		toTraderID := trader2.ID

		if i%2 == 1 {
			fromTraderID = trader2.ID
			toTraderID = trader1.ID
		}

		go func() {
			_, err := store.PaymentTx(context.Background(), PaymentTxParams{
				FromTraderID: fromTraderID,
				ToTraderID:   toTraderID,
				Amount:       amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// verify the final updated results
	updatedTrader1, err := store.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	updatedTrader2, err := store.GetTrader(context.Background(), trader2.ID)
	require.NoError(t, err)

	require.Equal(t, trader1.Balance, updatedTrader1.Balance)
	require.Equal(t, trader2.Balance, updatedTrader2.Balance)
}
