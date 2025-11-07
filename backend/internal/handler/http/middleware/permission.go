package middleware

import (
	"log"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/platform/i18n"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"

	"github.com/gofiber/fiber/v2"
)

func PermissionMiddleware(permService service.IPermissionService, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.Query("lang", "tr")

		user, ok := c.Locals("user").(*auth.AuthUser)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
				Success: false, Message: "Unauthorized",
			})
		}

		allowed, err := permService.CheckPermission(c.Context(), user.UserID, resource, action)
		if err != nil {
			log.Printf("Permission check error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
				Success: false,
				Message: i18n.Get(lang, "database_error"),
				Error:   &web.AppError{Code: "DATABASE_ERROR", Status: 500},
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(web.ApiResponse{
				Success: false,
				Message: i18n.Get(lang, "permission_denied"),
				Error:   &web.AppError{Code: "PERMISSION_DENIED", Status: 403},
			})
		}

		return c.Next()
	}
}
