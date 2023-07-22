package gapi

import (
	"context"
	"database/sql"

	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/pb"
	"github.com/YuanData/aquatrade/util"
	"github.com/YuanData/aquatrade/vldn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginMember(ctx context.Context, req *pb.LoginMemberRequest) (*pb.LoginMemberResponse, error) {
	violations := validateLoginMemberRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	member, err := server.store.GetMember(ctx, req.GetMembername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "member not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find member")
	}

	err = util.VerifyPassword(req.Password, member.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenGenerator.CreateToken(
		member.Membername,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenGenerator.CreateToken(
		member.Membername,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Membername:   member.Membername,
		RefreshToken: refreshToken,
		MemberAgent:  mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginMemberResponse{
		Member:                convertMember(member),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}

func validateLoginMemberRequest(req *pb.LoginMemberRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := vldn.ValidateMembername(req.GetMembername()); err != nil {
		violations = append(violations, fieldViolation("membername", err))
	}

	if err := vldn.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
