package service

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/repository"
)

// IUnitOfWorkFactory defines a factory for creating IUnitOfWork instances.
type IUnitOfWorkFactory interface {
	New(ctx context.Context) repository.IUnitOfWork
}

// unitOfWorkFactory is the concrete implementation of IUnitOfWorkFactory.
type unitOfWorkFactory struct {
	db *gorm.DB
}

// NewUnitOfWorkFactory creates a new unit of work factory.
func NewUnitOfWorkFactory(db *gorm.DB) IUnitOfWorkFactory {
	return &unitOfWorkFactory{db: db}
}

// New creates a new unit of work instance with a new transaction.
func (f *unitOfWorkFactory) New(ctx context.Context) repository.IUnitOfWork {
	return repository.NewUnitOfWork(f.db, ctx)
}
