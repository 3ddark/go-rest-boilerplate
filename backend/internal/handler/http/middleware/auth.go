package middleware

import (
	"strings"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/platform/web"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return web.Unauthorized(c, "Authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return web.Unauthorized(c, "Invalid authorization header format")
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return web.Unauthorized(c, "Invalid or expired token")
	}

	claims, ok := token.Claims.(*auth.JWTClaims)
	if !ok {
		return web.Unauthorized(c, "Invalid token claims")
	}

	c.Locals("user", &auth.AuthUser{
		UserID: claims.UserID,
		Email:  claims.Email,
	})

	return c.Next()
}
