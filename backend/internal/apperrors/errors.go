package apperrors

import "errors"

// Standart uygulama hataları
var (
	ErrNotFound             = errors.New("kayıt bulunamadı")
	ErrInvalidCredentials   = errors.New("geçersiz e-posta veya şifre")
	ErrEmailExists          = errors.New("e-posta adresi zaten kullanımda")
	ErrValidation           = errors.New("doğrulama hatası")
	ErrForbidden            = errors.New("yetkiniz yok")
	ErrUnauthorized         = errors.New("kimlik doğrulanmadı")
	ErrInternalServer       = errors.New("sunucu hatası")
	ErrInvalidRequest       = errors.New("geçersiz istek")
	Err2FASetupNotCompleted = errors.New("2FA kurulumu tamamlanmamış")
	ErrInvalid2FACode       = errors.New("geçersiz 2FA kodu")
)
