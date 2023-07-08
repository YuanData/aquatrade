package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoGenerator struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoGenerator(symmetricKey string) (Generator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d characters", chacha20poly1305.KeySize)
	}

	generator := &PasetoGenerator{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return generator, nil
}

func (generator *PasetoGenerator) CreateToken(membername string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(membername, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := generator.paseto.Encrypt(generator.symmetricKey, payload, nil)
	return token, payload, err
}

func (generator *PasetoGenerator) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := generator.paseto.Decrypt(token, generator.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
