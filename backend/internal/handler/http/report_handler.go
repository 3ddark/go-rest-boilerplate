package http

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/platform/i18n"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"
)

type ReportHandler struct {
	reportService service.IReportService
}

func NewReportHandler(reportService service.IReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) RequestReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	lang := c.Locals("lang").(string)

	// Basit bir DTO yerine şimdilik map kullanıyoruz
	var reqBody map[string]interface{}
	if err := c.BodyParser(&reqBody); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	reportType := "monthly_user_registrations" // Bu normalde request'ten gelir
	report, err := h.reportService.RequestReport(ctx, reportType, reqBody)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, i18n.Get(lang, "database_error"))
	}

	return web.Success(c, fiber.StatusAccepted, report, "Report generation started")
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	lang := c.Locals("lang").(string)

	id, err := c.ParamsInt("id")
	if err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, i18n.Get(lang, "invalid_request"))
	}

	report, err := h.reportService.GetReportStatus(ctx, id)
	if err != nil {
		return web.NotFound(c, i18n.Get(lang, "report_not_found"))
	}

	return web.Success(c, fiber.StatusOK, report)
}
