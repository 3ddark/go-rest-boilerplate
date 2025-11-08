package http

import (
	"github.com/gofiber/fiber/v2"
	"ths-erp.com/internal/domain"
	"ths-erp.com/internal/platform/web"
	"ths-erp.com/internal/service"
)

// LanguageHandler handles HTTP requests for languages.
type LanguageHandler struct {
	languageService service.ILanguageService
}

// NewLanguageHandler creates a new LanguageHandler instance.
func NewLanguageHandler(languageService service.ILanguageService) *LanguageHandler {
	return &LanguageHandler{languageService: languageService}
}

// GetAll handles the GET /api/v1/languages request.
func (h *LanguageHandler) GetAll(c *fiber.Ctx) error {
	// Default language is 'en'
	lang := c.Query("lang", "en")

	// Query parametrelerinden sayfalama bilgilerini al
	pagination := &domain.Pagination{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("pageSize", 10),
		SortBy:    c.Query("sortBy", "id"),
		SortOrder: c.Query("sortOrder", "asc"),
	}

	languages, pagination, err := h.languageService.GetActiveLanguages(c.Context(), lang, pagination)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Paginated(c, languages, pagination)
}
