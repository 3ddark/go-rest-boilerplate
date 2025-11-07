package domain

import "github.com/lib/pq"

// User, veritabanındaki 'users' tablosunu temsil eden ana modeldir.
type User struct {
	BaseEntity
	Name                   string         `json:"name" gorm:"column:name"`
	Email                  string         `json:"email" gorm:"column:email;unique"`
	PasswordHash           string         `json:"-" gorm:"column:password_hash"`
	TwoFactorEnabled       bool           `json:"twoFactorEnabled" gorm:"column:two_factor_enabled;default:false"`
	TwoFactorSecret        string         `json:"-" gorm:"column:two_factor_secret"`
	TwoFactorRecoveryCodes pq.StringArray `json:"-" gorm:"column:two_factor_recovery_codes;type:text[]"`
}

// UserPermission, bir kullanıcının belirli bir kaynak üzerindeki yetkilerini tanımlar.
type UserPermission struct {
	BaseEntity
	UserID     int    `json:"userId" gorm:"column:user_id"`
	Resource   string `json:"resource" gorm:"column:resource"`
	CanAdd     bool   `json:"canAdd" gorm:"column:can_add"`
	CanUpdate  bool   `json:"canUpdate" gorm:"column:can_update"`
	CanDelete  bool   `json:"canDelete" gorm:"column:can_delete"`
	CanSelect  bool   `json:"canSelect" gorm:"column:can_select"`
	CanSpecial bool   `json:"canSpecial" gorm:"column:can_special"`
}
