package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId uint) (*string, error) {
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTPayload{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	TOKEN_SECRET := os.Getenv("TOKEN_SECRET")
	token, err := claim.SignedString([]byte(TOKEN_SECRET))
	if err != nil {
		return nil, err
	}

	return &token, nil
}
