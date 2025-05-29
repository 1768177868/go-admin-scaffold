package v1

import (
	"net/http"

	"app/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// GetLocales returns a list of available locales
func GetLocales(c *gin.Context) {
	i18nInst := i18n.GetInstance()
	locales := i18nInst.GetLocales()

	c.JSON(http.StatusOK, gin.H{
		"locales": locales,
		"default": i18nInst.GetDefaultLocale(),
	})
}

// GetTranslations returns all translations for a specific locale
func GetTranslations(c *gin.Context) {
	locale := c.Query("locale")
	if locale == "" {
		locale = i18n.GetInstance().GetDefaultLocale()
	}

	translations := i18n.GetInstance().GetTranslations(locale)
	if translations == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Locale not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"locale":       locale,
		"translations": translations,
	})
}
