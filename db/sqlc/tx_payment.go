package db

import "context"

type PaymentTxParams struct {
	FromTraderID int64 `json:"from_account_id"`
	ToTraderID   int64 `json:"to_account_id"`
	Amount       int64 `json:"amount"`
}

type PaymentTxResult struct {
	Payment    Payment `json:"payment"`
	FromTrader Trader  `json:"from_account"`
	ToTrader   Trader  `json:"to_account"`
	FromRecord Record  `json:"from_record"`
	ToRecord   Record  `json:"to_record"`
}

func (store *Store) PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error) {
	var result PaymentTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Payment, err = q.CreatePayment(ctx, CreatePaymentParams{
			FromTraderID: arg.FromTraderID,
			ToTraderID:   arg.ToTraderID,
			Amount:       arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromRecord, err = q.CreateRecord(ctx, CreateRecordParams{
			TraderID: arg.FromTraderID,
			Amount:   -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToRecord, err = q.CreateRecord(ctx, CreateRecordParams{
			TraderID: arg.ToTraderID,
			Amount:   arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
