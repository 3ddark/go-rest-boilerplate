package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/platform/web"
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

	// Query parametrelerinden sayfalama bilgilerini al
	pagination := &domain.Pagination{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("pageSize", 10),
		SortBy:    c.Query("sortBy", "id"),
		SortOrder: c.Query("sortOrder", "asc"),
	}

	countries, pagination, err := h.countryService.GetAll(c.Context(), lang, pagination)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Paginated(c, countries, pagination)
}

// GetByCode handles the GET /api/v1/countries/:code request.
func (h *CountryHandler) GetByCode(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	code := c.Params("code")

	country, err := h.countryService.GetByCode(c.Context(), code, lang)
	if err != nil {
		return web.NotFound(c, "Country not found")
	}

	return web.Success(c, fiber.StatusOK, country)
}

// Create handles the POST /api/v1/countries request.
func (h *CountryHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCountryRequest
	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return web.ValidationError(c, err)
	}

	if err := h.countryService.Create(c.Context(), req); err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Success(c, fiber.StatusCreated, nil, "Country created successfully")
}

// Update handles the PUT /api/v1/countries/:code request.
func (h *CountryHandler) Update(c *fiber.Ctx) error {
	code := c.Params("code")
	var req dto.UpdateCountryRequest
	if err := c.BodyParser(&req); err != nil {
		return web.CustomError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return web.ValidationError(c, err)
	}

	if err := h.countryService.Update(c.Context(), code, req); err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Success(c, fiber.StatusOK, nil, "Country updated successfully")
}

// Delete handles the DELETE /api/v1/countries/:code request.
func (h *CountryHandler) Delete(c *fiber.Ctx) error {
	code := c.Params("code")
	if err := h.countryService.Delete(c.Context(), code); err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Success(c, fiber.StatusNoContent, nil, "Country deleted successfully")
}
