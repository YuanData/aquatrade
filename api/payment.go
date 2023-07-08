package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/token"

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

	fromTrader, valid := server.validTrader(ctx, req.FromTraderID, req.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromTrader.Holder != authPayload.Membername {
		err := errors.New("from trader doesn't belong to the authenticated member")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validTrader(ctx, req.ToTraderID, req.Currency)
	if !valid {
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

func (server *Server) validTrader(ctx *gin.Context, traderID int64, currency string) (db.Trader, bool) {
	trader, err := server.store.GetTrader(ctx, traderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return trader, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return trader, false
	}

	if trader.Currency != currency {
		err := fmt.Errorf("trader [%d] currency mismatch: %s vs %s", trader.ID, trader.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return trader, false
	}

	return trader, true
}
