package token

import (
	"testing"
	"time"

	"github.com/YuanData/aquatrade/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoGenerator(t *testing.T) {
	generator, err := NewPasetoGenerator(util.RandomString(32))
	require.NoError(t, err)

	membername := util.RandomHolder()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := generator.CreateToken(membername, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = generator.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, membername, payload.Membername)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	generator, err := NewPasetoGenerator(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := generator.CreateToken(util.RandomHolder(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = generator.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
