package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// IUnitOfWork defines the interface for a unit of work.
// It provides access to repositories and methods to commit or rollback the transaction.
type IUnitOfWork interface {
	UserRepository() IUserRepository
	ReportRepository() IReportRepository
	PermissionRepository() IPermissionRepository
	CountryRepository() ICountryRepository
	LanguageRepository() ILanguageRepository
	UnitRepository() IUnitRepository
	Commit() error
	Rollback()
}

// unitOfWork is the concrete implementation of IUnitOfWork.
type unitOfWork struct {
	db *gorm.DB
	tx *gorm.DB
}

// NewUnitOfWork creates a new unit of work instance.
// It begins a new transaction.
func NewUnitOfWork(db *gorm.DB, ctx context.Context) IUnitOfWork {
	tx := db.WithContext(ctx).Begin()
	return &unitOfWork{
		db: db,
		tx: tx,
	}
}

// UserRepository returns a user repository that uses the transaction.
func (u *unitOfWork) UserRepository() IUserRepository {
	return NewUserRepository(u.tx)
}

// ReportRepository returns a report repository that uses the transaction.
func (u *unitOfWork) ReportRepository() IReportRepository {
	return NewReportRepository(u.tx)
}

// PermissionRepository returns a permission repository that uses the transaction.
func (u *unitOfWork) PermissionRepository() IPermissionRepository {
	return NewPermissionRepository(u.tx)
}

// CountryRepository returns a country repository that uses the transaction.
func (u *unitOfWork) CountryRepository() ICountryRepository {
	return NewCountryRepository(u.tx)
}

// LanguageRepository returns a language repository that uses the transaction.
func (u *unitOfWork) LanguageRepository() ILanguageRepository {
	return NewLanguageRepository(u.tx)
}

// UnitRepository returns a unit repository that uses the transaction.
func (u *unitOfWork) UnitRepository() IUnitRepository {
	return NewUnitRepository(u.tx)
}

// Commit commits the transaction.
func (u *unitOfWork) Commit() error {
	if err := u.tx.Commit().Error; err != nil {
		// Rollback on commit error
		u.tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the transaction.
func (u *unitOfWork) Rollback() {
	u.tx.Rollback()
}
