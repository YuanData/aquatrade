package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createTraderRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTrader(ctx *gin.Context) {
	var req createTraderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateTraderParams{
		Holder:   authPayload.Membername,
		Currency: req.Currency,
		Balance:  0,
	}

	trader, err := server.store.CreateTrader(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trader)
}

type getTraderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTrader(ctx *gin.Context) {
	var req getTraderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trader, err := server.store.GetTrader(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if trader.Holder != authPayload.Membername {
		err := errors.New("trader doesn't belong to the authenticated member")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, trader)
}

type listTraderRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTraders(ctx *gin.Context) {
	var req listTraderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListTradersParams{
		Holder: authPayload.Membername,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	traders, err := server.store.ListTraders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, traders)
}
