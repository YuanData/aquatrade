package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/YuanData/aquatrade/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromTraderID int64  `json:"from_trader_id" binding:"required,min=1"`
	ToTraderID   int64  `json:"to_trader_id" binding:"required,min=1"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Currency     string `json:"currency" binding:"required,currency"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validTrader(ctx, req.FromTraderID, req.Currency) {
		return
	}

	if !server.validTrader(ctx, req.ToTraderID, req.Currency) {
		return
	}

	arg := db.PaymentTxParams{
		FromTraderID: req.FromTraderID,
		ToTraderID:   req.ToTraderID,
		Amount:       req.Amount,
	}

	result, err := server.store.PaymentTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validTrader(ctx *gin.Context, traderID int64, currency string) bool {
	trader, err := server.store.GetTrader(ctx, traderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if trader.Currency != currency {
		err := fmt.Errorf("trader [%d] currency mismatch: %s vs %s", trader.ID, trader.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
