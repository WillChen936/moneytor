package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func NewPayload(userID int64, duration time.Duration) *Payload {
	return &Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
