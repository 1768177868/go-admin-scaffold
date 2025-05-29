package i18n

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	defaultLocale = "en"
	instance      *I18n
	once          sync.Once
)

// I18n represents the internationalization instance
type I18n struct {
	sync.RWMutex
	defaultLocale string
	translations  map[string]map[string]interface{}
}

// Config represents i18n configuration
type Config struct {
	DefaultLocale    string   `mapstructure:"default_locale"`
	LoadPath         string   `mapstructure:"load_path"`
	AvailableLocales []string `mapstructure:"available_locales"`
}

// New creates a new I18n instance
func New(config *Config) *I18n {
	once.Do(func() {
		instance = &I18n{
			defaultLocale: config.DefaultLocale,
			translations:  make(map[string]map[string]interface{}),
		}
		if err := instance.loadTranslations(config.LoadPath); err != nil {
			panic(fmt.Sprintf("Failed to load translations: %v", err))
		}
	})
	return instance
}

// GetInstance returns the singleton instance of I18n
func GetInstance() *I18n {
	if instance == nil {
		panic("I18n not initialized")
	}
	return instance
}

// T translates a key into the specified locale
func (i *I18n) T(locale, key string, args ...interface{}) string {
	i.RLock()
	defer i.RUnlock()

	// If locale doesn't exist, fall back to default locale
	if _, ok := i.translations[locale]; !ok {
		locale = i.defaultLocale
	}

	// Split the key by dots to traverse nested translations
	parts := strings.Split(key, ".")
	value := i.getNestedValue(i.translations[locale], parts)

	// If value not found in specified locale, try default locale
	if value == "" && locale != i.defaultLocale {
		value = i.getNestedValue(i.translations[i.defaultLocale], parts)
	}

	// If still not found, return the key
	if value == "" {
		return key
	}

	// If there are arguments, format the translation
	if len(args) > 0 {
		return fmt.Sprintf(value, args...)
	}

	return value
}

// getNestedValue retrieves a nested value from a map using dot notation
func (i *I18n) getNestedValue(data map[string]interface{}, keys []string) string {
	if len(keys) == 0 {
		return ""
	}

	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key should return a string value
			if str, ok := current[key].(string); ok {
				return str
			}
			return ""
		}

		// For intermediate keys, get the nested map
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}

	return ""
}

// GetTranslations returns all translations for a specific locale
func (i *I18n) GetTranslations(locale string) map[string]interface{} {
	i.RLock()
	defer i.RUnlock()

	if translations, ok := i.translations[locale]; ok {
		return translations
	}
	return nil
}

// loadTranslations loads all translation files from the specified directory
func (i *I18n) loadTranslations(path string) error {
	files, err := filepath.Glob(filepath.Join(path, "*.yml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		locale := strings.TrimSuffix(filepath.Base(file), ".yml")

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read translation file %s: %v", file, err)
		}

		var translations map[string]interface{}
		if err := yaml.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse translation file %s: %v", file, err)
		}

		i.translations[locale] = translations
	}

	return nil
}

// GetLocales returns a list of available locales
func (i *I18n) GetLocales() []string {
	i.RLock()
	defer i.RUnlock()

	locales := make([]string, 0, len(i.translations))
	for locale := range i.translations {
		locales = append(locales, locale)
	}
	return locales
}

// GetDefaultLocale returns the default locale
func (i *I18n) GetDefaultLocale() string {
	return i.defaultLocale
}
