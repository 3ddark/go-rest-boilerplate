package dto

import "ths-erp.com/internal/domain"

// BaseResponse - Tüm response DTO'ları bundan türer
type BaseResponse struct {
	ID int `json:"id"`
}

func (b *BaseResponse) GetID() int {
	return b.ID
}

// LoginRequest - Kullanıcı girişi için kullanılan DTO.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse - Başarılı giriş sonrası dönen DTO.
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// CreateUserRequest - Yeni kullanıcı oluşturmak için kullanılan DTO.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest - Kullanıcı bilgilerini güncellemek için kullanılan DTO.
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse - API'den kullanıcı bilgisi dönerken kullanılan DTO.
// Bu DTO, domain.User modelindeki hassas bilgileri (örn: PasswordHash) dışarıya sızdırmaz.
type UserResponse struct {
	BaseResponse
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Derleme zamanında UserResponse'un domain.IResponse arayüzünü uyguladığını kontrol eder.
var _ domain.IResponse = (*UserResponse)(nil)
