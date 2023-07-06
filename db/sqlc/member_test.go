package db

import (
	"context"
	"testing"
	"time"

	"github.com/YuanData/aquatrade/util"
	"github.com/stretchr/testify/require"
)

func createRandomMember(t *testing.T) Member {
	arg := CreateMemberParams{
		Membername:     util.RandomHolder(),
		HashedPassword: "password",
		FullName:       util.RandomHolder(),
		Email:          util.RandomEmail(),
	}

	member, err := testQueries.CreateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, member)

	require.Equal(t, arg.Membername, member.Membername)
	require.Equal(t, arg.HashedPassword, member.HashedPassword)
	require.Equal(t, arg.FullName, member.FullName)
	require.Equal(t, arg.Email, member.Email)
	require.True(t, member.PasswordChangedAt.IsZero())
	require.NotZero(t, member.CreatedAt)

	return member
}

func TestCreateMember(t *testing.T) {
	createRandomMember(t)
}

func TestGetMember(t *testing.T) {
	member1 := createRandomMember(t)
	member2, err := testQueries.GetMember(context.Background(), member1.Membername)
	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.Membername, member2.Membername)
	require.Equal(t, member1.HashedPassword, member2.HashedPassword)
	require.Equal(t, member1.FullName, member2.FullName)
	require.Equal(t, member1.Email, member2.Email)
	require.WithinDuration(t, member1.PasswordChangedAt, member2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, member1.CreatedAt, member2.CreatedAt, time.Second)
}
