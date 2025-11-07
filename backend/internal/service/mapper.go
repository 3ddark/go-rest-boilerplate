package service

import (
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
)

// IMapper - Domain ve DTO nesneleri arasında dönüşüm yapmak için genel arayüz.
type IMapper[T domain.IEntity, R domain.IResponse] interface {
	ToResponse(entity T) R
	ToResponseList(entities []T) []R
	ToEntity(req domain.IRequest) T
	ToEntityUpdate(req domain.IRequest) T
}

// UserMapper - User için özel mapper implementasyonu.
type UserMapper struct{}

// NewUserMapper, UserMapper için bir kurucu fonksiyondur.
func NewUserMapper() IMapper[*domain.User, *dto.UserResponse] {
	return &UserMapper{}
}

func (m *UserMapper) ToResponse(user *domain.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		BaseResponse: dto.BaseResponse{ID: user.ID},
		Name:         user.Name,
		Email:        user.Email,
	}
}

func (m *UserMapper) ToResponseList(users []*domain.User) []*dto.UserResponse {
	responses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, m.ToResponse(user))
	}
	return responses
}

func (m *UserMapper) ToEntity(req domain.IRequest) *domain.User {
	createReq, ok := req.(*dto.CreateUserRequest)
	if !ok {
		return nil
	}
	// PasswordHash burada atanmaz, service katmanında hash'lendikten sonra atanır.
	return &domain.User{
		Name:  createReq.Name,
		Email: createReq.Email,
	}
}

func (m *UserMapper) ToEntityUpdate(req domain.IRequest) *domain.User {
	updateReq, ok := req.(*dto.UpdateUserRequest)
	if !ok {
		return nil
	}
	return &domain.User{
		Name:  updateReq.Name,
		Email: updateReq.Email,
	}
}
