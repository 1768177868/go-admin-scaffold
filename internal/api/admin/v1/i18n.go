package v1

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"app/pkg/response"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// GetLocales handles the request to get supported locales
// @Summary Get supported locales
// @Description Get list of supported locales
// @Tags i18n
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /admin/v1/i18n/locales [get]
func GetLocales(c *gin.Context) {
	locales := []string{"zh", "en"}
	response.Success(c, locales)
}

// GetTranslations handles the request to get translation data
// @Summary Get translations
// @Description Get translation data for specified locale
// @Tags i18n
// @Accept json
// @Produce json
// @Param locale query string true "Locale (e.g., zh, en)"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/v1/i18n/translations [get]
func GetTranslations(c *gin.Context) {
	locale := c.Query("locale")
	if locale == "" {
		response.ParamError(c, "locale parameter is required")
		return
	}

	// Validate locale
	if locale != "zh" && locale != "en" {
		response.ParamError(c, "unsupported locale")
		return
	}

	// Read translation file
	filePath := filepath.Join("locales", fmt.Sprintf("%s.yml", locale))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		response.Error(c, response.CodeNotFound, "translation file not found")
		return
	}

	// Parse YAML
	var translations map[string]interface{}
	if err := yaml.Unmarshal(data, &translations); err != nil {
		response.Error(c, response.CodeServerError, "failed to parse translation file")
		return
	}

	response.Success(c, translations)
}

// GetSupportedLocales handles the request to get supported locales
// @Summary Get supported locales
// @Description Get list of supported locales
// @Tags i18n
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /admin/v1/i18n/locales [get]
func GetSupportedLocales(c *gin.Context) {
	locales := []string{"zh", "en"}
	response.Success(c, locales)
}
