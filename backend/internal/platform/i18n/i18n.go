package i18n

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init() {
	bundle = i18n.NewBundle(language.Turkish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("internal/platform/i18n/locales/tr.json")
	bundle.LoadMessageFile("internal/platform/i18n/locales/en.json")
}

func Get(lang, key string) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})
}

func Middleware(c *fiber.Ctx) error {
	lang := c.Query("lang")
	if lang == "" {
		lang = c.Get("Accept-Language")
	}
	if lang == "" {
		lang = "tr"
	}
	c.Locals("lang", lang)
	return c.Next()
}
