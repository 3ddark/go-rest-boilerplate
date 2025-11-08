package http

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/apperrors"
	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/i18n"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"
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

// handleError, servis katmanından gelen hataları uygun HTTP yanıtlarına dönüştürür.
func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	lang := c.Locals("lang").(string)
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return web.NotFound(c, i18n.Get(lang, "user_not_found"))
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		return web.Unauthorized(c, i18n.Get(lang, "invalid_credentials"))
	case errors.Is(err, apperrors.ErrValidation):
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	case errors.Is(err, apperrors.ErrEmailExists):
		return web.CustomError(c, fiber.StatusConflict, i18n.Get(lang, "email_exists"))
	case errors.Is(err, apperrors.ErrInvalid2FACode):
		return web.Unauthorized(c, i18n.Get(lang, "invalid_2fa_code"))
	default:
		log.Printf("Unhandled error in UserHandler: %v", err)
		return web.CustomError(c, fiber.StatusInternalServerError, i18n.Get(lang, "internal_server_error"))
	}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	user, err := h.userService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return h.handleError(c, err)
	}

	if user.TwoFactorEnabled {
		return web.Success(c, fiber.StatusOK, fiber.Map{"userId": user.ID, "twoFactorRequired": true}, "2FA required")
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return web.CustomError(c, fiber.StatusInternalServerError, "Could not process login request")
	}

	userResponse := h.userMapper.ToResponse(user)
	loginResponse := dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	return web.Success(c, fiber.StatusOK, loginResponse, "Login successful")
}

func (h *UserHandler) Login2FA(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.Login2FARequest

	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	valid, err := h.userService.Verify2FA(ctx, req.UserID, req.Code)
	if err != nil || !valid {
		return h.handleError(c, err)
	}

	user, err := h.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return h.handleError(c, err)
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return web.CustomError(c, fiber.StatusInternalServerError, "Could not process login request")
	}

	loginResponse := dto.LoginResponse{
		Token: token,
		User:  user,
	}

	return web.Success(c, fiber.StatusOK, loginResponse, "Login successful")
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, user, i18n.Get(lang, "users_retrieved"))
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return h.handleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, users, i18n.Get(lang, "users_retrieved"))
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, user, i18n.Get(lang, "user_created"))
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	user, err := h.userService.UpdateUser(ctx, id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, user, i18n.Get(lang, "user_updated"))
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), handlerTimeout)
	defer cancel()

	lang := c.Locals("lang").(string)
	id, err := c.ParamsInt("id")
	if err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, nil, i18n.Get(lang, "user_deleted"))
}