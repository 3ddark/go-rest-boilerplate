package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/platform/metrics"
)

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	start := time.Now()
	var user domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	duration := time.Since(start).Seconds()
	metrics.M.DbQueryDuration.WithLabelValues("select", "users").Observe(duration)

	if result.Error == gorm.ErrRecordNotFound {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "not_found").Inc()
		return nil, result.Error
	}
	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "success").Inc()
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	start := time.Now()
	var user domain.User

	result := r.db.WithContext(ctx).First(&user, id)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("select", "users").Observe(duration)

	if result.Error == gorm.ErrRecordNotFound {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "not_found").Inc()
		return nil, result.Error
	}
	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "success").Inc()
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	start := time.Now()
	var users []domain.User

	result := r.db.WithContext(ctx).Find(&users)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("select", "users").Observe(duration)

	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	metrics.M.DbQueriesTotal.WithLabelValues("select", "users", "success").Inc()
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	start := time.Now()
	result := r.db.WithContext(ctx).Create(user)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("insert", "users").Observe(duration)

	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("insert", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	metrics.M.DbQueriesTotal.WithLabelValues("insert", "users", "success").Inc()
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id int, user *domain.User) (*domain.User, error) {
	start := time.Now()
	// Sadece belirtilen alanları güncellemek için Updates kullanılır.
	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(user)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("update", "users").Observe(duration)

	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("update", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		metrics.M.DbQueriesTotal.WithLabelValues("update", "users", "not_found").Inc()
		return nil, gorm.ErrRecordNotFound
	}

	// Güncellenmiş veriyi geri döndürmek için tekrar sorgu yapalım.
	var updatedUser domain.User
	r.db.WithContext(ctx).First(&updatedUser, id)

	metrics.M.DbQueriesTotal.WithLabelValues("update", "users", "success").Inc()
	return &updatedUser, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	result := r.db.WithContext(ctx).Delete(&domain.User{}, id)
	duration := time.Since(start).Seconds()

	metrics.M.DbQueryDuration.WithLabelValues("delete", "users").Observe(duration)

	if result.Error != nil {
		metrics.M.DbQueriesTotal.WithLabelValues("delete", "users", "error").Inc()
		metrics.M.DatabaseErrorsTotal.Inc()
		return result.Error
	}

	if result.RowsAffected == 0 {
		metrics.M.DbQueriesTotal.WithLabelValues("delete", "users", "not_found").Inc()
		return gorm.ErrRecordNotFound
	}

	metrics.M.DbQueriesTotal.WithLabelValues("delete", "users", "success").Inc()
	return nil
}
