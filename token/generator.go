package token

import (
	"time"
)

type Generator interface {
	CreateToken(membername string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
