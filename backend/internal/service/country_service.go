package service

import (
	"context"
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
	GetAll(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]dto.CountryResponse, *domain.Pagination, error)
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

func (s *CountryService) GetAll(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]dto.CountryResponse, *domain.Pagination, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	allowedSortBy := map[string]bool{
		"id":   true,
		"code": true,
	}
	if !allowedSortBy[pagination.SortBy] {
		pagination.SortBy = "id"
		pagination.SortOrder = "asc"
	}

	countries, pagination, err := uow.CountryRepository().FindAll(ctx, languageCode, pagination)
	if err != nil {
		return nil, nil, err
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

	return response, pagination, nil
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
