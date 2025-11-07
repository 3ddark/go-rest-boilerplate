package domain

// IEntity - Tüm domain modellerinin base interface'i
type IEntity interface {
	GetID() int
}

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
