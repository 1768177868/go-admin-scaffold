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
	translations  map[string]map[string]string
}

// Config represents i18n configuration
type Config struct {
	DefaultLocale string `mapstructure:"default_locale"`
	LoadPath      string `mapstructure:"load_path"`
}

// New creates a new I18n instance
func New(config *Config) *I18n {
	once.Do(func() {
		instance = &I18n{
			defaultLocale: config.DefaultLocale,
			translations:  make(map[string]map[string]string),
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

	// Get translation
	translation, ok := i.translations[locale][key]
	if !ok {
		// If translation not found in specified locale, try default locale
		if locale != i.defaultLocale {
			if t, ok := i.translations[i.defaultLocale][key]; ok {
				translation = t
			} else {
				return key
			}
		} else {
			return key
		}
	}

	// If there are arguments, format the translation
	if len(args) > 0 {
		return fmt.Sprintf(translation, args...)
	}

	return translation
}

// GetTranslations returns all translations for a specific locale
func (i *I18n) GetTranslations(locale string) map[string]string {
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

		var translations map[string]string
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
