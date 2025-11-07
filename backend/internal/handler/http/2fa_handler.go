
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
		return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	resp, err := h.userService.Setup2FA(ctx, user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false,
			Message: "Could not setup 2FA",
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true,
		Message: "2FA setup initiated",
		Data:    resp,
	})
}

func (h *UserHandler) Enable2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	user, err := auth.GetUserFromContext(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	var req dto.Enable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false,
			Message: "Invalid request",
		})
	}

	recoveryCodes, err := h.userService.Enable2FA(ctx, user.UserID, req.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true,
		Message: "2FA enabled successfully",
		Data:    recoveryCodes,
	})
}

func (h *UserHandler) Disable2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	user, err := auth.GetUserFromContext(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	if err := h.userService.Disable2FA(ctx, user.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false,
			Message: "Could not disable 2FA",
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true,
		Message: "2FA disabled successfully",
	})
}
