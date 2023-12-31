// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	Membername        string    `json:"membername"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Payment struct {
	ID           int64 `json:"id"`
	FromTraderID int64 `json:"from_trader_id"`
	ToTraderID   int64 `json:"to_trader_id"`
	// most be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Record struct {
	ID       int64 `json:"id"`
	TraderID int64 `json:"trader_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Membername   string    `json:"membername"`
	RefreshToken string    `json:"refresh_token"`
	MemberAgent  string    `json:"member_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Trader struct {
	ID        int64     `json:"id"`
	Holder    string    `json:"holder"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
