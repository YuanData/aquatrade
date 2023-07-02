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
	}
}
