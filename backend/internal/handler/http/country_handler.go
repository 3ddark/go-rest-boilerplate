
package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/service"
)

// CountryHandler handles HTTP requests for countries.
type CountryHandler struct {
	countryService service.ICountryService
	validate       *validator.Validate
}

// NewCountryHandler creates a new CountryHandler instance.
func NewCountryHandler(countryService service.ICountryService) *CountryHandler {
	return &CountryHandler{
		countryService: countryService,
		validate:       validator.New(),
	}
}

// GetAll handles the GET /api/v1/countries request.
func (h *CountryHandler) GetAll(c *fiber.Ctx) error {
	// Default language is 'en'
	lang := c.Query("lang", "en")

	countries, err := h.countryService.GetAll(c.Context(), lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(countries)
}

// GetByCode handles the GET /api/v1/countries/:code request.
func (h *CountryHandler) GetByCode(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	code := c.Params("code")

	country, err := h.countryService.GetByCode(c.Context(), code, lang)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Country not found",
		})
	}

	return c.JSON(country)
}

// Create handles the POST /api/v1/countries request.
func (h *CountryHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCountryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.countryService.Create(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Country created successfully"})
}

// Update handles the PUT /api/v1/countries/:code request.
func (h *CountryHandler) Update(c *fiber.Ctx) error {
	code := c.Params("code")
	var req dto.UpdateCountryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.countryService.Update(c.Context(), code, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Country updated successfully"})
}

// Delete handles the DELETE /api/v1/countries/:code request.
func (h *CountryHandler) Delete(c *fiber.Ctx) error {
	code := c.Params("code")
	if err := h.countryService.Delete(c.Context(), code); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Country deleted successfully"})
}

