package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"ths-erp.com/internal/platform/logger" // Kendi modül adınızla
)

type contextKey string

const loggerKey contextKey = "logger"

// AttachLogger, her isteğe özel bir logger oluşturur ve context'e ekler.
func AttachLogger(c *fiber.Ctx) error {
	reqID := uuid.New().String()
	c.Set("X-Request-ID", reqID)

	// İstekle ilgili alanları içeren bir alt-logger (sub-logger) oluştur
	reqLogger := logger.L.With().
		Str("request_id", reqID).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Logger()

	// Logger'ı context'e ekle
	// c.UserContext() zaten bir context.Context döndürür.
	ctx := context.WithValue(c.UserContext(), loggerKey, &reqLogger)
	c.SetUserContext(ctx)

	start := time.Now()

	// Bir sonraki middleware veya handler'a devam et
	err := c.Next()

	// İstek bittiğinde son log'u at
	reqLogger.Info().
		Int("status", c.Response().StatusCode()).
		Dur("latency_ms", time.Since(start)).
		Msg("Request finished")

	return err
}

// GetLogger, context'ten isteğe özel logger'ı alır.
// Eğer context'ten logger alınamazsa, global fallback logger'ı döner.
// DÜZELTME: ctx context.T -> ctx context.Context
func GetLogger(ctx context.Context) *zerolog.Logger {
	// Nil context durumuna karşı koruma (panic'i önler)
	if ctx == nil {
		return &logger.L
	}

	if l, ok := ctx.Value(loggerKey).(*zerolog.Logger); ok && l != nil {
		return l
	}

	// Fallback olarak global logger'ı döndür
	return &logger.L
}
