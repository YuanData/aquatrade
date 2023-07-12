package gapi

import (
	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertMember(member db.Member) *pb.Member {
	return &pb.Member{
		Membername:        member.Membername,
		FullName:          member.FullName,
		Email:             member.Email,
		PasswordChangedAt: timestamppb.New(member.PasswordChangedAt),
		CreatedAt:         timestamppb.New(member.CreatedAt),
	}
}
