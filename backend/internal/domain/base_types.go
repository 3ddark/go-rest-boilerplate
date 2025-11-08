package domain

import (
	"fmt"
)

// Pagination - Sayfalama bilgilerini tutar
type Pagination struct {
	Page         int    `json:"page"`
	PageSize     int    `json:"pageSize"`
	SortBy       string `json:"sortBy"`
	SortOrder    string `json:"sortOrder"` // "asc" or "desc"
	TotalRecords int64  `json:"totalRecords"`
	TotalPages   int    `json:"totalPages"`
}

// GetOffset - GORM için offset değerini hesaplar
func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit - GORM için limit değerini döndürür
func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// GetSort - GORM için sıralama string'ini oluşturur
func (p *Pagination) GetSort() string {
	if p.SortBy == "" {
		return "id asc"
	}
	if p.SortOrder == "" {
		p.SortOrder = "asc"
	}
	return fmt.Sprintf("%s %s", p.SortBy, p.SortOrder)
}

// IEntity - Tüm domain modellerinin base interface'i
type IEntity interface {
	GetID() int
}

// ... (rest of the file)

// BaseEntity - Tüm domain modelleri bundan türer
type BaseEntity struct {
	ID int `json:"id" gorm:"primaryKey"`
}

func (b *BaseEntity) GetID() int {
	return b.ID
}

// IRequest - Tüm request DTO'larının base interface'i
type IRequest interface{}

// IResponse - Tüm response DTO'larının base interface'i
type IResponse interface {
	GetID() int
}
