package main

import (
	"log"
	"os"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/config"
	"ths-erp.com/internal/handler/graphql"
	"ths-erp.com/internal/handler/http"
	"ths-erp.com/internal/platform/database"
	"ths-erp.com/internal/platform/database/migration"
	"ths-erp.com/internal/platform/i18n"
	"ths-erp.com/internal/platform/logger"
	"ths-erp.com/internal/platform/metrics"
	"ths-erp.com/internal/platform/queue"
	"ths-erp.com/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logger.Init()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	auth.SetJWTSecret(cfg.JWTSecret)

	// Initialize prometheus metrics platforms
	i18n.Init()
	metrics.Init()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✓ Database connected")

	// Run migrations
	migration.Migrate(db)

	// Connect to RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}
	rabbitClient, err := queue.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	err = rabbitClient.DeclareAndBindQueue("app_exchange", "welcome_emails_queue", "user.welcome_email")
	if err != nil {
		log.Fatalf("Could not declare queue/exchange: %v", err)
	}

	// Mappers
	userMapper := service.NewUserMapper()

	// Cache
	// appCache := cache.NewInMemoryCache(5 * time.Minute)

	// Unit of Work Factory
	uowFactory := service.NewUnitOfWorkFactory(db)

	// Services
	permService := service.NewPermissionService(uowFactory)
	userService := service.NewUserService(uowFactory, userMapper, rabbitClient)

	// Setup server
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Routes
	http.SetupRoutes(app, db, permService, rabbitClient)
	graphql.SetupHandler(app, userService, permService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✓ Server running on port %s", port)
	log.Printf("✓ Metrics available at http://localhost:%s/metrics", port)
	log.Printf("✓ GraphQL UI available at http://localhost:%s/graphiql", port)
	log.Fatal(app.Listen(":" + port))
}
