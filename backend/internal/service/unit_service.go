package service

import (
	"context"

	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
)

type IUnitService interface {
	GetAllUnits(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]dto.UnitDTO, *domain.Pagination, error)
}

type unitService struct {
	uowFactory IUnitOfWorkFactory
}

func NewUnitService(uowFactory IUnitOfWorkFactory) IUnitService {
	return &unitService{uowFactory}
}

func (s *unitService) GetAllUnits(ctx context.Context, languageCode string, pagination *domain.Pagination) ([]dto.UnitDTO, *domain.Pagination, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	// SQL Injection'ı önlemek için sıralama kolonunu beyaz listede kontrol et
	allowedSortBy := map[string]bool{
		"id":   true,
		"code": true,
	}
	if !allowedSortBy[pagination.SortBy] {
		pagination.SortBy = "id" // Geçersizse varsayılana dön
		pagination.SortOrder = "asc"
	}

	units, pagination, err := uow.UnitRepository().FindAll(ctx, languageCode, pagination)
	if err != nil {
		return nil, nil, err
	}

	var unitDTOs []dto.UnitDTO
	for _, unit := range units {
		var translatedName string
		if len(unit.Translations) > 0 {
			translatedName = unit.Translations[0].Name
		} else {
			translatedName = "N/A"
		}
		unitDTOs = append(unitDTOs, dto.UnitDTO{
			Code: unit.Code,
			Name: translatedName,
		})
	}

	return unitDTOs, pagination, nil
}
