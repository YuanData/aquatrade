package db

import "context"

type PaymentTxParams struct {
	FromTraderID int64 `json:"from_trader_id"`
	ToTraderID   int64 `json:"to_trader_id"`
	Amount       int64 `json:"amount"`
}

type PaymentTxResult struct {
	Payment    Payment `json:"payment"`
	FromTrader Trader  `json:"from_trader"`
	ToTrader   Trader  `json:"to_trader"`
	FromRecord Record  `json:"from_record"`
	ToRecord   Record  `json:"to_record"`
}

func (store *SQLStore) PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error) {
	var result PaymentTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createPaymentParams := CreatePaymentParams(arg)
		result.Payment, err = q.CreatePayment(ctx, createPaymentParams)
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

		if arg.FromTraderID < arg.ToTraderID {
			result.FromTrader, result.ToTrader, err = transferAmount(ctx, q, arg.FromTraderID, -arg.Amount, arg.ToTraderID, arg.Amount)
		} else {
			result.ToTrader, result.FromTrader, err = transferAmount(ctx, q, arg.ToTraderID, arg.Amount, arg.FromTraderID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func transferAmount(ctx context.Context, q *Queries, traderID1 int64, amount1 int64, traderID2 int64, amount2 int64,
) (trader1 Trader, trader2 Trader, err error) {
	trader1, err = q.AddTraderBalance(ctx, AddTraderBalanceParams{
		ID:     traderID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	trader2, err = q.AddTraderBalance(ctx, AddTraderBalanceParams{
		ID:     traderID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
