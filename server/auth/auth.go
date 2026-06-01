package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
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

func ValidateJWT(token string) (*JWTPayload, error) {
	TOKEN_SECRET := os.Getenv("TOKEN_SECRET")
	claims := &JWTPayload{}
	token0, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_SECRET), nil
	})
	if err != nil || !token0.Valid {
		return nil, err
	}

	return claims, nil
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "No token provided")
		}
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		c.Set("userID", claims.UserID)

		return next(c)
	}
}
