package gapi

import (
	"context"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/pb"
	"github.com/YuanData/aquatrade/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.CreateMemberResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateMemberParams{
		Membername:     req.GetMembername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	member, err := server.store.CreateMember(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "membername already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create member: %s", err)
	}

	rsp := &pb.CreateMemberResponse{
		Member: convertMember(member),
	}
	return rsp, nil
}
