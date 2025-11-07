package repository

import (
	"context"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

type ILanguageRepository interface {
	GetActiveLanguages(ctx context.Context, translationLanguageCode string) ([]domain.Language, error)
}

type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) ILanguageRepository {
	return &LanguageRepository{db: db}
}

func (r *LanguageRepository) GetActiveLanguages(ctx context.Context, translationLanguageCode string) ([]domain.Language, error) {
	var languages []domain.Language
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Preload("Translations", "translation_language_code = ?", translationLanguageCode).
		Find(&languages).Error

	if err != nil {
		return nil, err
	}
	return languages, nil
}
