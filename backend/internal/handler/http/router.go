package http

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
	"ths-erp.com/internal/handler/http/middleware"
	"ths-erp.com/internal/platform/cache"
	"ths-erp.com/internal/platform/queue"
	"ths-erp.com/internal/service"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, permService service.IPermissionService, queueClient *queue.RabbitMQClient, redisClient *redis.Client) {
	uowFactory := service.NewUnitOfWorkFactory(db)
	appCache := cache.NewRedisCache(redisClient)

	// Initialize services
	userService := service.NewUserService(uowFactory, &service.UserMapper{}, queueClient)
	countryService := service.NewCountryService(uowFactory, appCache)
	languageService := service.NewLanguageService(uowFactory, appCache)
	unitService := service.NewUnitService(uowFactory)
	reportService := service.NewReportService(uowFactory, queueClient)

	// Initialize handlers
	userHandler := NewUserHandler(userService, permService, &service.UserMapper{})
	countryHandler := NewCountryHandler(countryService)
	languageHandler := NewLanguageHandler(languageService)
	unitHandler := NewUnitHandler(unitService)
	reportHandler := NewReportHandler(reportService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes
	v1.Post("/login", userHandler.Login)
	v1.Get("/countries", countryHandler.GetAll)
	v1.Get("/languages", languageHandler.GetAll)
	v1.Get("/units", unitHandler.GetUnits)

	v1.Use(middleware.AuthMiddleware)

	userHandler.Setup2FARoutes(v1)

	countryRoutes := v1.Group("/countries")
	countryRoutes.Get("/:code", countryHandler.GetByCode)
	countryRoutes.Post("/", middleware.PermissionMiddleware(permService, "country", "write"), countryHandler.Create)
	countryRoutes.Put("/:code", middleware.PermissionMiddleware(permService, "country", "write"), countryHandler.Update)
	countryRoutes.Delete("/:code", middleware.PermissionMiddleware(permService, "country", "delete"), countryHandler.Delete)

	userRoutes := v1.Group("/users")
	userRoutes.Get("/", middleware.PermissionMiddleware(permService, "user", "read"), userHandler.GetAll)
	userRoutes.Get("/:id", middleware.PermissionMiddleware(permService, "user", "read"), userHandler.Get)
	userRoutes.Post("/", middleware.PermissionMiddleware(permService, "user", "write"), userHandler.Create)
	userRoutes.Put("/:id", middleware.PermissionMiddleware(permService, "user", "write"), userHandler.Update)
	userRoutes.Delete("/:id", middleware.PermissionMiddleware(permService, "user", "delete"), userHandler.Delete)

	reportRoutes := v1.Group("/reports")
	reportRoutes.Post("/", reportHandler.RequestReport)
	reportRoutes.Get("/:id", reportHandler.GetReport)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
