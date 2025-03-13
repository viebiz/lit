package i18n

import (
	"context"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/viebiz/lit/monitoring"
	"golang.org/x/text/language"
)

const (
	defaultLang             = "en"
	defaultBundleFileFormat = "json"
)

type BundleConfig struct {
	DefaultLang            string
	SourcePath             string
	BundleFileFormat       string
	ExtraBundleFileSupport map[string]UnmarshalFunc
}

func Init(ctx context.Context, cfg BundleConfig) MessageBundle {
	// Create a new bundleMessage
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc(defaultBundleFileFormat, json.Unmarshal)

	// Register custom unmarshal function
	for format, unmarshalFunc := range cfg.ExtraBundleFileSupport {
		b.RegisterUnmarshalFunc(format, unmarshalFunc)
	}

	if cfg.DefaultLang == "" {
		cfg.DefaultLang = defaultLang
	}

	if cfg.BundleFileFormat == "" {
		cfg.BundleFileFormat = defaultBundleFileFormat
	}

	return &bundleMessage{
		underlyingBundle: b,
		DefaultLang:      cfg.DefaultLang,
		SourcePath:       cfg.SourcePath,
		BundleFileFormat: cfg.BundleFileFormat,
		monitor:          monitoring.FromContext(ctx),
		localizeMap:      make(map[string]Localizable),
	}
}
