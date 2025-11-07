package service

import (
	"context"

	"ths-erp.com/internal/dto"
)

type IUnitService interface {
	GetAllUnits(ctx context.Context, languageCode string) ([]dto.UnitDTO, error)
}

type unitService struct {
	uowFactory IUnitOfWorkFactory
}

func NewUnitService(uowFactory IUnitOfWorkFactory) IUnitService {
	return &unitService{uowFactory}
}

func (s *unitService) GetAllUnits(ctx context.Context, languageCode string) ([]dto.UnitDTO, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	units, err := uow.UnitRepository().FindAll(ctx, languageCode)
	if err != nil {
		return nil, err
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

	return unitDTOs, nil
}
