package api

import (
	"database/sql"
	"net/http"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createTraderRequest struct {
	Holder   string `json:"holder" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTrader(ctx *gin.Context) {
	var req createTraderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTraderParams{
		Holder:   req.Holder,
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

	arg := db.ListTradersParams{
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
