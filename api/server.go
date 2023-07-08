package api

import (
	"fmt"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/token"
	"github.com/YuanData/aquatrade/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config         util.Config
	store          db.Store
	router         *gin.Engine
	tokenGenerator token.Generator
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewJWTGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token : %w", err)
	}
	server := &Server{
		config:         config,
		store:          store,
		tokenGenerator: tokenGenerator,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenGenerator))
	authRoutes.POST("/traders", server.createTrader)
	authRoutes.GET("/traders/:id", server.getTrader)
	authRoutes.GET("/traders", server.listTraders)
	authRoutes.POST("/payments", server.createPayment)

	router.POST("/members", server.createMember)
	router.POST("/members/login", server.loginMember)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
