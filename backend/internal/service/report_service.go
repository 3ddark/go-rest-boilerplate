package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/platform/queue"
)

// GenerateReportJob, rapor oluşturma görevi için kuyruğa atılacak veriyi tanımlar.
type GenerateReportJob struct {
	ReportID int `json:"report_id"`
}

type IReportService interface {
	RequestReport(ctx context.Context, reportType string, payload map[string]interface{}) (*domain.Report, error)
	GetReportStatus(ctx context.Context, id int) (*domain.Report, error)
	ProcessReport(ctx context.Context, reportID int) error // Bu Worker tarafından çağrılacak
}

type ReportService struct {
	uowFactory  IUnitOfWorkFactory
	queueClient *queue.RabbitMQClient
}

func NewReportService(uowFactory IUnitOfWorkFactory, queueClient *queue.RabbitMQClient) IReportService {
	return &ReportService{
		uowFactory:  uowFactory,
		queueClient: queueClient,
	}
}

// RequestReport API tarafından çağrılır.
func (s *ReportService) RequestReport(ctx context.Context, reportType string, payload map[string]interface{}) (*domain.Report, error) {
	payloadBytes, _ := json.Marshal(payload)

	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	// 1. Veritabanına rapor kaydı oluştur (status: pending)
	report := &domain.Report{
		Type:    reportType,
		Status:  domain.ReportStatusPending,
		Payload: string(payloadBytes),
	}
	createdReport, err := uow.ReportRepository().Create(ctx, report)
	if err != nil {
		return nil, err
	}

	// 2. Değişiklikleri commit et
	if err := uow.Commit(); err != nil {
		return nil, err
	}

	// 3. RabbitMQ'ya görevi yayınla (Commit'ten sonra)
	job := GenerateReportJob{ReportID: createdReport.ID}
	jobPayload, _ := json.Marshal(job)

	err = s.queueClient.Publish(ctx, "app_exchange", "report.generate", jobPayload)
	if err != nil {
		log.Printf("ERROR: Could not publish report job for report %d: %v", createdReport.ID, err)
		// Burada telafi edici bir işlem düşünülebilir. Örneğin, rapor durumunu 'failed' olarak güncellemek.
	}

	return createdReport, nil
}

// GetReportStatus API tarafından çağrılır.
func (s *ReportService) GetReportStatus(ctx context.Context, id int) (*domain.Report, error) {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()
	return uow.ReportRepository().GetByID(ctx, id)
}

// ProcessReport Worker tarafından çağrılır.
func (s *ReportService) ProcessReport(ctx context.Context, reportID int) error {
	uow := s.uowFactory.New(ctx)
	defer uow.Rollback()

	reportRepo := uow.ReportRepository()

	// 1. Raporu DB'den al
	report, err := reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}

	// 2. Durumu 'processing' yap
	report.Status = domain.ReportStatusProcessing
	if err := reportRepo.Update(ctx, report); err != nil {
		return err
	}

	// 3. Ağır işi yap: Raporu oluştur
	userRepo := uow.UserRepository()
	allUsers, err := userRepo.FindAll(ctx) // Bu normalde filtrelenmiş bir sorgu olmalı
	if err != nil {
		report.Status = domain.ReportStatusFailed
		report.Error = err.Error()
		_ = reportRepo.Update(ctx, report) // Hata durumunda durumu güncelle, ana hatayı döndür
		return uow.Commit()                // Hata durumunu kaydetmek için commit et
	}

	resultData := map[string]interface{}{
		"report_type":  report.Type,
		"total_users":  len(allUsers),
		"generated_at": time.Now(),
	}
	resultBytes, _ := json.Marshal(resultData)

	// 4. Sonucu ve durumu DB'ye kaydet
	report.Status = domain.ReportStatusCompleted
	report.Result = resultBytes
	report.Error = ""
	if err := reportRepo.Update(ctx, report); err != nil {
		return err
	}

	// 5. Tüm değişiklikleri commit et
	return uow.Commit()
}
