package service

import (
	"context"
)

type IPermissionService interface {
	CheckPermission(ctx context.Context, userID int, resource, action string) (bool, error)
}

type PermissionService struct {
	uowFactory IUnitOfWorkFactory
}

func NewPermissionService(uowFactory IUnitOfWorkFactory) IPermissionService {
	return &PermissionService{uowFactory: uowFactory}
}

func (s *PermissionService) CheckPermission(ctx context.Context, userID int, resource, action string) (bool, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback() // Read-only operation

	// Bu katman şimdilik direkt repository'i çağırıyor,
	// ileride cache'leme gibi ek iş mantıkları buraya eklenebilir.
	return uow.PermissionRepository().CheckPermission(ctx, userID, resource, action)
}
