package worker

import (
	"context"
	"encoding/json"
	"time"

	"ths-erp.com/internal/platform/logger"
	"ths-erp.com/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

// JobConsumer, kuyruktan gelen görevleri işler.
type JobConsumer struct {
	channel       *amqp.Channel
	userService   service.IUserService
	reportService service.IReportService
}

// NewJobConsumer, JobConsumer için bir kurucu fonksiyondur.
func NewJobConsumer(channel *amqp.Channel, userService service.IUserService, reportService service.IReportService) *JobConsumer {
	return &JobConsumer{
		channel:       channel,
		userService:   userService,
		reportService: reportService,
	}
}

// StartConsumers, projedeki tüm görev kuyruklarını dinlemeye başlar.
// Her kuyruk kendi goroutine'inde çalışır, böylece birbirlerini bloklamazlar.
func (c *JobConsumer) StartConsumers() {
	go c.consume("welcome_emails_queue", c.handleWelcomeEmail)
	go c.consume("reports_queue", c.handleGenerateReport)

	logger.L.Info().Msg("All consumers started. Waiting for messages...")

	// Ana goroutine'in sonlanmasını engellemek için sonsuz bir bekleme döngüsü.
	// Worker'ın sürekli çalışmasını sağlar.
	select {}
}

// consume, belirtilen bir kuyruğu dinleyen ve gelen her mesaj için
// verilen handler fonksiyonunu çağıran jenerik bir metottur.
func (c *JobConsumer) consume(queueName string, handler func(d amqp.Delivery)) {
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer tag
		false,     // auto-ack: false olmalı, çünkü işlemi biz onaylayacağız.
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.L.Fatal().Err(err).Str("queue", queueName).Msg("Failed to register a consumer")
	}

	for delivery := range msgs {
		handler(delivery)
	}
}

// handleWelcomeEmail, 'welcome_emails_queue' kuyruğundan gelen mesajları işler.
func (c *JobConsumer) handleWelcomeEmail(d amqp.Delivery) {
	l := logger.L.With().Str("job_type", "welcome_email").Logger()
	l.Info().RawJSON("body", d.Body).Msg("Received a welcome email job")

	var job service.WelcomeEmailJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		l.Error().Err(err).Msg("Failed to unmarshal message. Rejecting.")
		d.Reject(false) // 'false' ile tekrar kuyruğa alınmaz, çünkü mesaj bozuk.
		return
	}

	// Asıl işi burada yapıyoruz (simülasyon)
	if err := c.sendWelcomeEmail(job); err != nil {
		l.Error().Err(err).Int("user_id", job.UserID).Msg("Failed to process job. Nacking.")
		// Nack (Negative Acknowledge) ile tekrar denemek üzere kuyruğa geri gönderebiliriz (requeue=true)
		// veya Dead Letter Exchange'e göndermek için (requeue=false) ayarlayabiliriz.
		d.Nack(false, false) // Şimdilik tekrar denemiyoruz.
	} else {
		// İş başarıyla bittiyse, mesajı kuyruktan siliyoruz (Acknowledge).
		d.Ack(false)
		l.Info().Int("user_id", job.UserID).Msg("Job processed successfully.")
	}
}

// handleGenerateReport, 'reports_queue' kuyruğundan gelen mesajları işler.
func (c *JobConsumer) handleGenerateReport(d amqp.Delivery) {
	l := logger.L.With().Str("job_type", "generate_report").Logger()
	l.Info().RawJSON("body", d.Body).Msg("Received a report generation job")

	var job service.GenerateReportJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		l.Error().Err(err).Msg("Failed to unmarshal message. Rejecting.")
		d.Reject(false)
		return
	}

	// Rapor oluşturma gibi uzun sürebilecek işlemler için timeout'lu bir context oluşturuyoruz.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Asıl işi ReportService'e devrediyoruz
	if err := c.reportService.ProcessReport(ctx, job.ReportID); err != nil {
		l.Error().Err(err).Int("report_id", job.ReportID).Msg("Failed to process report job. Nacking.")
		d.Nack(false, false)
	} else {
		d.Ack(false)
		l.Info().Int("report_id", job.ReportID).Msg("Report job processed successfully.")
	}
}

// sendWelcomeEmail, e-posta gönderme işini simüle eden bir yardımcı fonksiyondur.
// Gerçekte bu mantık IUserService içinde olmalıdır.
func (c *JobConsumer) sendWelcomeEmail(job service.WelcomeEmailJob) error {
	logger.L.Info().Str("email", job.Email).Int("user_id", job.UserID).Msg("Sending welcome email...")

	// Gerçek bir e-posta gönderme servisi çağrısını simüle edelim.
	time.Sleep(2 * time.Second)

	logger.L.Info().Str("email", job.Email).Msg("Email successfully sent.")
	return nil
}
