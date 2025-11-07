package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/cache"
)

const (
	countriesCacheKeyPrefix = "countries:"
	countriesCacheDuration  = 24 * time.Hour
)

type ICountryService interface {
	GetAll(ctx context.Context, languageCode string) ([]dto.CountryResponse, error)
	GetByCode(ctx context.Context, code string, languageCode string) (*dto.CountryResponse, error)
	Create(ctx context.Context, req dto.CreateCountryRequest) error
	Update(ctx context.Context, code string, req dto.UpdateCountryRequest) error
	Delete(ctx context.Context, code string) error
}

type CountryService struct {
	uowFactory IUnitOfWorkFactory
	cache      cache.ICache
}

func NewCountryService(uowFactory IUnitOfWorkFactory, cache cache.ICache) ICountryService {
	return &CountryService{
		uowFactory: uowFactory,
		cache:      cache,
	}
}

func (s *CountryService) GetAll(ctx context.Context, languageCode string) ([]dto.CountryResponse, error) {
	cacheKey := fmt.Sprintf("%s%s", countriesCacheKeyPrefix, languageCode)

	if cachedData, found := s.cache.Get(cacheKey); found {
		if data, ok := cachedData.([]byte); ok {
			var countries []dto.CountryResponse
			if err := json.Unmarshal(data, &countries); err == nil {
				return countries, nil
			}
		}
	}

	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	countries, err := uow.CountryRepository().FindAll(ctx, languageCode)
	if err != nil {
		return nil, err
	}

	response := make([]dto.CountryResponse, 0, len(countries))
	for _, country := range countries {
		name := country.Code
		if len(country.Translations) > 0 {
			name = country.Translations[0].Name
		}
		response = append(response, dto.CountryResponse{
			Code: country.Code,
			Name: name,
		})
	}

	if data, err := json.Marshal(response); err == nil {
		s.cache.Set(cacheKey, data, countriesCacheDuration)
	}

	return response, nil
}

func (s *CountryService) GetByCode(ctx context.Context, code string, languageCode string) (*dto.CountryResponse, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	country, err := uow.CountryRepository().FindByCode(ctx, code, languageCode)
	if err != nil {
		return nil, err
	}

	name := country.Code
	if len(country.Translations) > 0 {
		name = country.Translations[0].Name
	}

	return &dto.CountryResponse{
		Code: country.Code,
		Name: name,
	}, nil
}

func (s *CountryService) Create(ctx context.Context, req dto.CreateCountryRequest) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	country := domain.Country{
		Code: req.Code,
		Translations: []domain.CountryTranslation{
			{
				LanguageCode: "tr",
				Name:         req.Name,
			},
		},
	}

	if err := uow.CountryRepository().Create(ctx, country); err != nil {
		return err
	}

	return uow.Commit()
}

func (s *CountryService) Update(ctx context.Context, code string, req dto.UpdateCountryRequest) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	country, err := uow.CountryRepository().FindByCode(ctx, code, "tr")
	if err != nil {
		return err
	}

	if len(country.Translations) > 0 {
		country.Translations[0].Name = req.Name
	} else {
		country.Translations = append(country.Translations, domain.CountryTranslation{
			LanguageCode: "tr",
			Name:         req.Name,
		})
	}

	if err := uow.CountryRepository().Update(ctx, *country); err != nil {
		return err
	}

	return uow.Commit()
}

func (s *CountryService) Delete(ctx context.Context, code string) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	if err := uow.CountryRepository().DeleteByCode(ctx, code); err != nil {
		return err
	}

	return uow.Commit()
}
