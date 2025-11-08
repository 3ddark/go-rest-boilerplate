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
			return web.Unauthorized(c)
		}

		allowed, err := permService.CheckPermission(c.Context(), user.UserID, resource, action)
		if err != nil {
			log.Printf("Permission check error: %v", err)
			return web.CustomError(c, fiber.StatusInternalServerError, i18n.Get(lang, "database_error"))
		}

		if !allowed {
			return web.Forbidden(c, i18n.Get(lang, "permission_denied"))
		}

		return c.Next()
	}
}
