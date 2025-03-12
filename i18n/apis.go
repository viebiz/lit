package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	defaultLangTag          = "en"
	defaultBundleFileFormat = "json"
)

type BundleConfig struct {
	DefaultLang            string
	ExtraBundleFileSupport map[string]UnmarshalFunc
}

func Init(cfg BundleConfig) Bundle {
	// Create a new bundle
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc(defaultBundleFileFormat, json.Unmarshal)

	// Register custom unmarshal function
	for format, unmarshalFunc := range cfg.ExtraBundleFileSupport {
		b.RegisterUnmarshalFunc(format, unmarshalFunc)
	}

	defaultLang := defaultLangTag
	if cfg.DefaultLang != "" {
		defaultLang = cfg.DefaultLang
	}

	return &bundle{
		i18nBundle:  b,
		DefaultLang: defaultLang,
		LocalizeMap: make(map[string]MessageLocalize),
	}
}
