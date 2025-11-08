package service

import (
	"context"
	"time"

	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/cache"
)

const (
	languagesCacheKeyPrefix = "languages:active:"
	languagesCacheDuration  = 24 * time.Hour
)

type ILanguageService interface {
	GetActiveLanguages(ctx context.Context, translationLanguageCode string, pagination *domain.Pagination) ([]dto.LanguageResponse, *domain.Pagination, error)
}

type LanguageService struct {
	uowFactory IUnitOfWorkFactory
	cache      cache.ICache
}

func NewLanguageService(uowFactory IUnitOfWorkFactory, cache cache.ICache) ILanguageService {
	return &LanguageService{
		uowFactory: uowFactory,
		cache:      cache,
	}
}

func (s *LanguageService) GetActiveLanguages(ctx context.Context, translationLanguageCode string, pagination *domain.Pagination) ([]dto.LanguageResponse, *domain.Pagination, error) {
	// Sayfalama nedeniyle cache'leme kaldırıldı.
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback() // Read-only operation

	// SQL Injection'ı önlemek için sıralama kolonunu beyaz listede kontrol et
	allowedSortBy := map[string]bool{
		"id":   true,
		"code": true,
	}
	if !allowedSortBy[pagination.SortBy] {
		pagination.SortBy = "id" // Geçersizse varsayılana dön
		pagination.SortOrder = "asc"
	}

	languages, pagination, err := uow.LanguageRepository().GetActiveLanguages(ctx, translationLanguageCode, pagination)
	if err != nil {
		return nil, nil, err
	}

	response := make([]dto.LanguageResponse, 0, len(languages))
	for _, lang := range languages {
		name := lang.Code // Fallback to code if no translation is found
		if len(lang.Translations) > 0 {
			name = lang.Translations[0].Name
		}
		response = append(response, dto.LanguageResponse{
			Code: lang.Code,
			Name: name,
		})
	}

	return response, pagination, nil
}
