package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image/png"
	"log"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"ths-erp.com/internal/apperrors"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/metrics"
	"ths-erp.com/internal/platform/queue"
)

type WelcomeEmailJob struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type IUserService interface {
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
	Setup2FA(ctx context.Context, userID int) (*dto.Setup2FAResponse, error)
	Enable2FA(ctx context.Context, userID int, code string) ([]string, error)
	Disable2FA(ctx context.Context, userID int) error
	Verify2FA(ctx context.Context, userID int, code string) (bool, error)
}

type UserService struct {
	uowFactory  IUnitOfWorkFactory
	mapper      IMapper[*domain.User, *dto.UserResponse]
	queueClient *queue.RabbitMQClient
}

func NewUserService(uowFactory IUnitOfWorkFactory, mapper IMapper[*domain.User, *dto.UserResponse], queueClient *queue.RabbitMQClient) IUserService {
	return &UserService{
		uowFactory:  uowFactory,
		mapper:      mapper,
		queueClient: queueClient,
	}
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	userRepo := uow.UserRepository()
	user, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*dto.UserResponse, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	user, err := uow.UserRepository().FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return s.mapper.ToResponse(user), nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	users, err := uow.UserRepository().FindAll(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		resp := s.mapper.ToResponse(&users[i])
		if resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		metrics.M.ValidationErrorsTotal.WithLabelValues("user", "missing_fields").Inc()
		return nil, apperrors.ErrValidation
	}

	if !govalidator.IsEmail(req.Email) {
		metrics.M.ValidationErrorsTotal.WithLabelValues("email", "invalid_format").Inc()
		return nil, apperrors.ErrValidation
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error while hashing password: %v", err)
		return nil, apperrors.ErrInternalServer
	}

	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	// Check if email exists
	_, err = uow.UserRepository().FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, apperrors.ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // A different database error occurred
	}

	user := s.mapper.ToEntity(req)
	user.PasswordHash = string(hashedPassword)

	createdUser, err := uow.UserRepository().Create(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := uow.Commit(); err != nil {
		return nil, err
	}

	job := WelcomeEmailJob{
		UserID: createdUser.ID,
		Email:  createdUser.Email,
		Name:   createdUser.Name,
	}
	jobPayload, err := json.Marshal(job)
	if err != nil {
		log.Printf("ERROR: Could not marshal welcome email job for user %d: %v", createdUser.ID, err)
	} else {
		err := s.queueClient.Publish(ctx, "app_exchange", "user.welcome_email", jobPayload)
		if err != nil {
			log.Printf("ERROR: Could not publish welcome email job for user %d: %v", createdUser.ID, err)
		} else {
			log.Printf("✓ Welcome email job published for user %d", createdUser.ID)
		}
	}

	return s.mapper.ToResponse(createdUser), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if req.Email != "" && !govalidator.IsEmail(req.Email) {
		metrics.M.ValidationErrorsTotal.WithLabelValues("email", "invalid_format").Inc()
		return nil, apperrors.ErrValidation
	}

	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	updateData := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	updatedUser, err := uow.UserRepository().Update(ctx, id, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if err := uow.Commit(); err != nil {
		return nil, err
	}

	return s.mapper.ToResponse(updatedUser), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	err := uow.UserRepository().Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return err
	}

	return uow.Commit()
}

func (s *UserService) Setup2FA(ctx context.Context, userID int) (*dto.Setup2FAResponse, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	userRepo := uow.UserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "THS-ERP",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	user.TwoFactorSecret = key.Secret()
	if _, err := userRepo.Update(ctx, user.ID, user); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if err := png.Encode(&buf, img); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	if err := uow.Commit(); err != nil {
		return nil, err
	}

	return &dto.Setup2FAResponse{
		Secret:   key.Secret(),
		QRCode:   "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		Recovery: nil,
	}, nil
}

func (s *UserService) Enable2FA(ctx context.Context, userID int, code string) ([]string, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	userRepo := uow.UserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if user.TwoFactorSecret == "" {
		return nil, apperrors.Err2FASetupNotCompleted
	}

	valid := totp.Validate(code, user.TwoFactorSecret)
	if !valid {
		return nil, apperrors.ErrInvalid2FACode
	}

	user.TwoFactorEnabled = true

	recoveryCodes := make([]string, 10)
	for i := 0; i < 10; i++ {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return nil, apperrors.ErrInternalServer
		}
		recoveryCodes[i] = strings.ToLower(base64.RawURLEncoding.EncodeToString(b))
	}
	user.TwoFactorRecoveryCodes = recoveryCodes

	if _, err := userRepo.Update(ctx, user.ID, user); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	if err := uow.Commit(); err != nil {
		return nil, err
	}

	return recoveryCodes, nil
}

func (s *UserService) Disable2FA(ctx context.Context, userID int) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	userRepo := uow.UserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return err
	}

	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""
	user.TwoFactorRecoveryCodes = nil

	if _, err := userRepo.Update(ctx, user.ID, user); err != nil {
		return apperrors.ErrInternalServer
	}

	return uow.Commit()
}

func (s *UserService) Verify2FA(ctx context.Context, userID int, code string) (bool, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	userRepo := uow.UserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, apperrors.ErrNotFound
		}
		return false, err
	}

	if !user.TwoFactorEnabled {
		return true, nil
	}

	if totp.Validate(code, user.TwoFactorSecret) {
		return true, nil
	}

	for i, recoveryCode := range user.TwoFactorRecoveryCodes {
		if code == recoveryCode {
			user.TwoFactorRecoveryCodes = append(user.TwoFactorRecoveryCodes[:i], user.TwoFactorRecoveryCodes[i+1:]...)
			if _, err := userRepo.Update(ctx, user.ID, user); err != nil {
				return false, apperrors.ErrInternalServer
			}
			if err := uow.Commit(); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, apperrors.ErrInvalid2FACode
}
