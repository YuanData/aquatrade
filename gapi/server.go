package gapi

import (
	"fmt"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/pb"
	"github.com/YuanData/aquatrade/token"
	"github.com/YuanData/aquatrade/util"
)

type Server struct {
	pb.UnimplementedAquaTradeServer
	config         util.Config
	store          db.Store
	tokenGenerator token.Generator
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:         config,
		store:          store,
		tokenGenerator: tokenGenerator,
	}

	return server, nil
}
