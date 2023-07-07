package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createMemberRequest struct {
	Membername string `json:"membername" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=5"`
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}

type memberResponse struct {
	Membername        string    `json:"membername"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newMemberResponse(member db.Member) memberResponse {
	return memberResponse{
		Membername:        member.Membername,
		FullName:          member.FullName,
		Email:             member.Email,
		PasswordChangedAt: member.PasswordChangedAt,
		CreatedAt:         member.CreatedAt,
	}
}

func (server *Server) createMember(ctx *gin.Context) {
	var req createMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateMemberParams{
		Membername:     req.Membername,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	member, err := server.store.CreateMember(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newMemberResponse(member)
	ctx.JSON(http.StatusOK, rsp)
}

type loginMemberRequest struct {
	Membername string `json:"membername" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6"`
}

type loginMemberResponse struct {
	AccessToken string         `json:"access_token"`
	Member      memberResponse `json:"member"`
}

func (server *Server) loginMember(ctx *gin.Context) {
	var req loginMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	member, err := server.store.GetMember(ctx, req.Membername)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(req.Password, member.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenGenerator.CreateToken(
		member.Membername,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginMemberResponse{
		AccessToken: accessToken,
		Member:      newMemberResponse(member),
	}
	ctx.JSON(http.StatusOK, rsp)
}
