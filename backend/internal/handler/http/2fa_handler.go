
package http

import (
	"context"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/web"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Setup2FARoutes(router fiber.Router) {
	router.Post("/2fa/setup", h.Setup2FA)
	router.Post("/2fa/enable", h.Enable2FA)
	router.Post("/2fa/disable", h.Disable2FA)
}

func (h *UserHandler) Setup2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	user, err := auth.GetUserFromContext(c.UserContext())
	if err != nil {
		return web.Unauthorized(c)
	}

	resp, err := h.userService.Setup2FA(ctx, user.UserID)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, "Could not setup 2FA")
	}

	return web.Success(c, fiber.StatusOK, resp, "2FA setup initiated")
}

func (h *UserHandler) Enable2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	user, err := auth.GetUserFromContext(c.UserContext())
	if err != nil {
		return web.Unauthorized(c)
	}

	var req dto.Enable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, "Invalid request")
	}

	recoveryCodes, err := h.userService.Enable2FA(ctx, user.UserID, req.Code)
	if err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, err.Error())
	}

	return web.Success(c, fiber.StatusOK, recoveryCodes, "2FA enabled successfully")
}

func (h *UserHandler) Disable2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	user, err := auth.GetUserFromContext(c.UserContext())
	if err != nil {
		return web.Unauthorized(c)
	}

	if err := h.userService.Disable2FA(ctx, user.UserID); err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, "Could not disable 2FA")
	}

	return web.Success(c, fiber.StatusOK, nil, "2FA disabled successfully")
}
