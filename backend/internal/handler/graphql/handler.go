package graphql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/graphql-go/graphql"
)

// SetupHandler, GraphQL endpoint'ini ve GraphiQL arayüzünü Fiber app'e ekler.
func SetupHandler(app *fiber.App, userService service.IUserService, permService service.IPermissionService) {
	schema, err := buildGraphQLSchema(userService, permService)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	// Gelen isteği parse etmek için kullanılacak struct.
	type GraphQLRequest struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	// GraphQL endpoint'i
	app.All("/graphql", authMiddleware, func(c *fiber.Ctx) error {
		// Middleware'den gelen kullanıcıyı al
		user := c.Locals("user")
		ctx := context.Background()
		if user != nil {
			// GraphQL context'ine kullanıcıyı ekle
			ctx = context.WithValue(ctx, auth.UserContextKey, user)
		}

		var req GraphQLRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			Context:        ctx,
		})

		return c.JSON(result)
	})

	// GraphiQL arayüzü
	app.Get("/graphiql", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(graphiqlHTML)
	})
}

// authMiddleware, GraphQL endpoint'i için özel bir JWT doğrulama middleware'idir.
// REST API'den farklı olarak, token olmasa bile devam eder, sadece context'e kullanıcı eklemez.
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Next() // Hatalı format, yine de devam et.
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.GetJWTSecret(), nil
	})

	// Token geçersizse veya süresi dolmuşsa, yine de devam et.
	// Resolver'lar context'te kullanıcı olup olmadığını kontrol ederek yetkilendirme yapar.
	if err != nil || !token.Valid {
		return c.Next()
	}

	claims, ok := token.Claims.(*auth.JWTClaims)
	if !ok {
		return c.Next()
	}

	c.Locals("user", &auth.AuthUser{
		UserID: claims.UserID,
		Email:  claims.Email,
	})

	return c.Next()
}
