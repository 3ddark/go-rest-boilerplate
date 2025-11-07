package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GetJWTSecret() []byte {
	return jwtSecret
}

type contextKey string

const UserContextKey contextKey = "user"

type AuthUser struct {
	UserID int
	Email  string
}

func GetUserFromContext(ctx context.Context) (*AuthUser, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	user, ok := ctx.Value(UserContextKey).(*AuthUser)
	if !ok || user == nil {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	return user, nil
}

func GenerateJWT(userID int, email string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
