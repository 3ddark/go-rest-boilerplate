package repository

import (
	"context"
	"time"

	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/platform/metrics"

	"gorm.io/gorm"
)

type IPermissionRepository interface {
	GetUserPermission(ctx context.Context, userID int, resource string) (*domain.UserPermission, error)
	CheckPermission(ctx context.Context, userID int, resource, action string) (bool, error)
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) IPermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetUserPermission(ctx context.Context, userID int, resource string) (*domain.UserPermission, error) {
	start := time.Now()
	var perm domain.UserPermission

	result := r.db.WithContext(ctx).Where("user_id = ? AND resource = ?", userID, resource).First(&perm)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("select", "user_permissions").Observe(duration)

	if result.Error == gorm.ErrRecordNotFound {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "user_permissions", "not_found").Inc()
		return nil, nil // Hata değil, sadece kayıt yok
	}
	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "user_permissions", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	metrics.M.DbQueriesTotal.WithLabelValues("select", "user_permissions", "success").Inc()
	return &perm, nil
}

func (r *PermissionRepository) CheckPermission(ctx context.Context, userID int, resource, action string) (bool, error) {
	perm, err := r.GetUserPermission(ctx, userID, resource)
	if err != nil {
		return false, err
	}

	if perm == nil {
		metrics.M.PermissionChecksTotal.WithLabelValues(resource, action, "false").Inc()
		return false, nil
	}

	allowed := false
	switch action {
	case "select":
		allowed = perm.CanSelect
	case "add":
		allowed = perm.CanAdd
	case "update":
		allowed = perm.CanUpdate
	case "delete":
		allowed = perm.CanDelete
	case "special":
		allowed = perm.CanSpecial
	}

	allowedStr := "false"
	if allowed {
		allowedStr = "true"
	}
	metrics.M.PermissionChecksTotal.WithLabelValues(resource, action, allowedStr).Inc()

	return allowed, nil
}
