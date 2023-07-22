package gapi

import (
	"context"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/pb"
	"github.com/YuanData/aquatrade/util"
	"github.com/YuanData/aquatrade/vldn"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.CreateMemberResponse, error) {
	violations := validateCreateMemberRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
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

func validateCreateMemberRequest(req *pb.CreateMemberRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := vldn.ValidateMembername(req.GetMembername()); err != nil {
		violations = append(violations, fieldViolation("membername", err))
	}

	if err := vldn.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := vldn.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := vldn.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
