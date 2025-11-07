package repository

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

type ICountryRepository interface {
	FindAll(ctx context.Context, languageCode string) ([]domain.Country, error)
	FindByCode(ctx context.Context, code string, languageCode string) (*domain.Country, error)
	Create(ctx context.Context, country domain.Country) error
	Update(ctx context.Context, country domain.Country) error
	Delete(ctx context.Context, id int64) error
	DeleteByCode(ctx context.Context, code string) error
}

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) ICountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) FindAll(ctx context.Context, languageCode string) ([]domain.Country, error) {
	var countries []domain.Country
	err := r.db.WithContext(ctx).Preload("Translations", "language_code = ?", languageCode).Find(&countries).Error
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *CountryRepository) FindByCode(ctx context.Context, code string, languageCode string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.WithContext(ctx).Preload("Translations", "language_code = ?", languageCode).First(&country, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) Create(ctx context.Context, country domain.Country) error {
	return r.db.WithContext(ctx).Create(&country).Error
}

func (r *CountryRepository) Update(ctx context.Context, country domain.Country) error {
	return r.db.WithContext(ctx).Save(&country).Error
}

func (r *CountryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Country{}, "id = ?", id).Error
}

func (r *CountryRepository) DeleteByCode(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Delete(&domain.Country{}, "code = ?", code).Error
}
