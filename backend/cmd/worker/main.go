package main

import (
	"log"
	"os"

	"ths-erp.com/internal/config"
	"ths-erp.com/internal/platform/database"
	"ths-erp.com/internal/platform/logger"
	"ths-erp.com/internal/platform/queue"
	"ths-erp.com/internal/service"
	"ths-erp.com/internal/worker"
)

func main() {
	logger.Init()
	logger.L.Info().Msg("Starting Worker process...")

	// 1. Konfigürasyonu Yükle
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Veritabanı Bağlantısı (Worker'ın da servislere ihtiyacı olabilir)
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✓ Database connected for worker")

	// 3. RabbitMQ Bağlantısı
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}
	rabbitClient, err := queue.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	// Tüketici olarak kuyruk ve exchange'in var olduğundan emin olalım
	err = rabbitClient.DeclareAndBindQueue("app_exchange", "welcome_emails_queue", "user.welcome_email")
	if err != nil {
		log.Fatalf("Could not declare queue/exchange: %v", err)
	}

	// 4. Bağımlılıkları Oluştur
	uowFactory := service.NewUnitOfWorkFactory(db)
	userMapper := service.NewUserMapper()

	// Worker'ın RabbitMQ'ya mesaj GÖNDERMESİNE gerek olmadığı için nil geçiyoruz.
	// Eğer worker başka bir görevi tetikleyecek olsaydı, client'ı buraya da geçerdik.
	userService := service.NewUserService(uowFactory, userMapper, nil)
	reportService := service.NewReportService(uowFactory, nil)

	// 5. Consumer'ı Başlat
	jobConsumer := worker.NewJobConsumer(rabbitClient.Channel, userService, reportService)
	jobConsumer.StartConsumers()
}
