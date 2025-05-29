package v1

import (
	"app/pkg/i18n"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetLocales returns a list of available locales
func GetLocales(c *gin.Context) {
	i18nInstance := i18n.GetInstance()
	response.Success(c, gin.H{
		"locales":        i18nInstance.GetLocales(),
		"default_locale": i18nInstance.GetDefaultLocale(),
	})
}

// GetTranslations returns all translations for a specific locale
func GetTranslations(c *gin.Context) {
	locale := c.Param("locale")
	i18nInstance := i18n.GetInstance()

	translations := i18nInstance.GetTranslations(locale)
	if translations == nil {
		response.NotFoundError(c)
		return
	}

	response.Success(c, gin.H{
		"locale":       locale,
		"translations": translations,
	})
}
