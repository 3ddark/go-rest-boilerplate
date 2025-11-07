package http

import (
	"context"
	"errors"
	"log"
	"time"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/i18n"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const handlerTimeout = 3 * time.Second

type UserHandler struct {
	userService service.IUserService
	userMapper  service.IMapper[*domain.User, *dto.UserResponse]
}

func NewUserHandler(userService service.IUserService, permService service.IPermissionService, mapper service.IMapper[*domain.User, *dto.UserResponse]) *UserHandler {
	return &UserHandler{
		userService: userService,
		userMapper:  mapper,
	}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false,
			Message: i18n.Get(lang, "invalid_request"),
		})
	}

	user, err := h.userService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	if user.TwoFactorEnabled {
		return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
			Success: true,
			Message: "2FA required",
			Data:    fiber.Map{"userId": user.ID, "twoFactorRequired": true},
		})
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false,
			Message: "Could not process login request",
		})
	}

	userResponse := h.userMapper.ToResponse(user)
	loginResponse := dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}

func (h *UserHandler) Login2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.Login2FARequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false,
			Message: i18n.Get(lang, "invalid_request"),
		})
	}

	valid, err := h.userService.Verify2FA(ctx, req.UserID, req.Code)
	if err != nil || !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Success: false,
			Message: "Invalid 2FA code",
		})
	}

	user, err := h.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false,
			Message: "Could not process login request",
		})
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false,
			Message: "Could not process login request",
		})
	}

	loginResponse := dto.LoginResponse{
		Token: token,
		User:  user,
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "invalid_request"),
		})
	}

	user, err := h.userService.GetUser(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "user_not_found"),
		})
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "database_error"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true, Message: i18n.Get(lang, "users_retrieved"), Data: user,
	})
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "database_error"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true, Message: i18n.Get(lang, "users_retrieved"), Data: users,
	})
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "invalid_request"),
		})
	}

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		// Service katmanından gelen validasyon hatalarını kontrol et
		if err.Error() == "invalid email format" {
			return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
				Success: false, Message: i18n.Get(lang, "invalid_email"),
			})
		}
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "database_error"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(web.ApiResponse{
		Success: true, Message: i18n.Get(lang, "user_created"), Data: user,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "invalid_request"),
		})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "invalid_request"),
		})
	}

	user, err := h.userService.UpdateUser(ctx, id, &req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "user_not_found"),
		})
	}
	if err != nil {
		if err.Error() == "invalid email format" {
			return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
				Success: false, Message: i18n.Get(lang, "invalid_email"),
			})
		}
		log.Printf("Error updating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "database_error"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true, Message: i18n.Get(lang, "user_updated"), Data: user,
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "invalid_request"),
		})
	}

	err = h.userService.DeleteUser(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "user_not_found"),
		})
	}
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Success: false, Message: i18n.Get(lang, "database_error"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Success: true, Message: i18n.Get(lang, "user_deleted"),
	})
}