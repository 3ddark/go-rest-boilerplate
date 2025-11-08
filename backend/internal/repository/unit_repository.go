package repository

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

type IUnitRepository interface {
	FindAll(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]domain.Unit, *domain.Pagination, error)
	FindByCode(ctx context.Context, code string, languageCode string) (*domain.Unit, error)
	Create(ctx context.Context, unit domain.Unit) error
	Update(ctx context.Context, unit domain.Unit) error
	Delete(ctx context.Context, id int64) error
	DeleteByCode(ctx context.Context, code string) error
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) IUnitRepository {
	return &unitRepository{db}
}

func (r *unitRepository) FindAll(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]domain.Unit, *domain.Pagination, error) {
	var units []domain.Unit
	var totalRecords int64

	// Toplam kayıt sayısını al
	if err := r.db.WithContext(ctx).Model(&domain.Unit{}).Count(&totalRecords).Error; err != nil {
		return nil, nil, err
	}

	pagination.TotalRecords = totalRecords
	pagination.TotalPages = int(totalRecords) / pagination.GetLimit()
	if int(totalRecords)%pagination.GetLimit() > 0 {
		pagination.TotalPages++
	}

	// Veriyi sayfalama ile al
	err := r.db.WithContext(ctx).
		Preload("Translations", "language_code = ?", languageCode).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order(pagination.GetSort()).
		Find(&units).Error

	if err != nil {
		return nil, nil, err
	}
	return units, pagination, nil
}

func (r *unitRepository) FindByCode(ctx context.Context, code string, languageCode string) (*domain.Unit, error) {
	var unit domain.Unit
	err := r.db.WithContext(ctx).Preload("Translations", "language_code = ?", languageCode).First(&unit, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Create(ctx context.Context, unit domain.Unit) error {
	return r.db.WithContext(ctx).Create(&unit).Error
}

func (r *unitRepository) Update(ctx context.Context, unit domain.Unit) error {
	return r.db.WithContext(ctx).Save(&unit).Error
}

func (r *unitRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Unit{}, "id = ?", id).Error
}

func (r *unitRepository) DeleteByCode(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Delete(&domain.Unit{}, "code = ?", code).Error
}
