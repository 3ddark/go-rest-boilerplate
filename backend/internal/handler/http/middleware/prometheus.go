package middleware

import (
	"fmt"
	"time"

	"ths-erp.com/internal/platform/metrics"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	metrics.M.HttpRequestsInFlight.Inc()
	defer metrics.M.HttpRequestsInFlight.Dec()

	err := c.Next()

	duration := time.Since(start).Seconds()
	endpoint := c.Path()
	method := c.Method()
	status := fmt.Sprintf("%d", c.Response().StatusCode())

	metrics.M.HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	metrics.M.HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()

	return err
}
