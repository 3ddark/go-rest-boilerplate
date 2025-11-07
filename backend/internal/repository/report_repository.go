package repository

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

type IReportRepository interface {
	Create(ctx context.Context, report *domain.Report) (*domain.Report, error)
	GetByID(ctx context.Context, id int) (*domain.Report, error)
	Update(ctx context.Context, report *domain.Report) error
}

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) IReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	if err := r.db.WithContext(ctx).Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportRepository) GetByID(ctx context.Context, id int) (*domain.Report, error) {
	var report domain.Report
	if err := r.db.WithContext(ctx).First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) Update(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}
