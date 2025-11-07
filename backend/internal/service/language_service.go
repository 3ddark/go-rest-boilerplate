package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/cache"
)

const (
	languagesCacheKeyPrefix = "languages:active:"
	languagesCacheDuration  = 24 * time.Hour
)

type ILanguageService interface {
	GetActiveLanguages(ctx context.Context, translationLanguageCode string) ([]dto.LanguageResponse, error)
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

func (s *LanguageService) GetActiveLanguages(ctx context.Context, translationLanguageCode string) ([]dto.LanguageResponse, error) {
	cacheKey := fmt.Sprintf("%s%s", languagesCacheKeyPrefix, translationLanguageCode)

	// 1. Try to get from cache
	if cachedData, found := s.cache.Get(cacheKey); found {
		if data, ok := cachedData.([]byte); ok {
			var languages []dto.LanguageResponse
			if err := json.Unmarshal(data, &languages); err == nil {
				return languages, nil
			}
		}
	}

	// 2. If not in cache, get from database
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback() // Read-only operation

	languages, err := uow.LanguageRepository().GetActiveLanguages(ctx, translationLanguageCode)
	if err != nil {
		return nil, err
	}

	// 3. Map to DTO
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

	// 4. Set to cache
	if data, err := json.Marshal(response); err == nil {
		s.cache.Set(cacheKey, data, languagesCacheDuration)
	}

	return response, nil
}
