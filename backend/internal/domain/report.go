package domain

import (
	"encoding/json"
	"time"
)

type ReportStatus string

const (
	ReportStatusPending    ReportStatus = "pending"
	ReportStatusProcessing ReportStatus = "processing"
	ReportStatusCompleted  ReportStatus = "completed"
	ReportStatusFailed     ReportStatus = "failed"
)

// Report, asenkron olarak oluşturulan bir raporu temsil eder.
type Report struct {
	BaseEntity
	Type      string          `json:"type" gorm:"column:type"`
	Status    ReportStatus    `json:"status" gorm:"column:status"`
	Payload   string          `json:"payload" gorm:"column:payload"` // Raporu oluşturmak için gereken parametreler (JSON)
	Result    json.RawMessage `json:"result" gorm:"column:result"`   // Raporun sonucu (JSON)
	Error     string          `json:"error,omitempty" gorm:"column:error"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}
