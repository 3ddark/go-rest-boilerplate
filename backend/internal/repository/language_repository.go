package repository

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

type ILanguageRepository interface {
	GetActiveLanguages(ctx context.Context, translationLanguageCode string, pagination *domain.Pagination) ([]domain.Language, *domain.Pagination, error)
}

type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) ILanguageRepository {
	return &LanguageRepository{db: db}
}

func (r *LanguageRepository) GetActiveLanguages(ctx context.Context, translationLanguageCode string, pagination *domain.Pagination) ([]domain.Language, *domain.Pagination, error) {
	var languages []domain.Language
	var totalRecords int64

	// Toplam aktif dil sayısını al
	query := r.db.WithContext(ctx).Model(&domain.Language{}).Where("is_active = ?", true)
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, nil, err
	}

	pagination.TotalRecords = totalRecords
	pagination.TotalPages = int(totalRecords) / pagination.GetLimit()
	if int(totalRecords)%pagination.GetLimit() > 0 {
		pagination.TotalPages++
	}

	// Aktif dilleri sayfalama ile al
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Preload("Translations", "translation_language_code = ?", translationLanguageCode).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order(pagination.GetSort()).
		Find(&languages).Error

	if err != nil {
		return nil, nil, err
	}
	return languages, pagination, nil
}
