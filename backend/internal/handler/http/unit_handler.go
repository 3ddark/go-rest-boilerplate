package http

import (
	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"
)

// UnitHandler handles HTTP requests for units.
type UnitHandler struct {
	unitService service.IUnitService
}

// NewUnitHandler creates a new instance of UnitHandler.
func NewUnitHandler(unitService service.IUnitService) *UnitHandler {
	return &UnitHandler{unitService}
}

// GetUnits returns a list of all units.
// @Summary Get all units
// @Description Get a list of all units of measurement with translations.
// @Tags Units
// @Accept  json
// @Produce  json
// @Param lang query string false "Language code for translations (e.g., en, tr)" default(en)
// @Success 200 {object} web.Response
// @Router /units [get]
func (h *UnitHandler) GetUnits(c *fiber.Ctx) error {
	languageCode := c.Query("lang", "en") // Default to 'en' if not provided

	units, err := h.unitService.GetAllUnits(c.Context(), languageCode)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, "Could not retrieve units")
	}

	return web.Success(c, fiber.StatusOK, units)
}
