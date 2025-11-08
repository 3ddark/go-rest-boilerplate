package http

import (
	"github.com/gofiber/fiber/v2"
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

	languages, err := h.languageService.GetActiveLanguages(c.Context(), lang)
	if err != nil {
		return web.CustomError(c, fiber.StatusInternalServerError, err.Error())
	}

	return web.Success(c, fiber.StatusOK, languages)
}
