package http

import (
	"context"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/i18n"
)

const (
	defaultSourcePath         = "resources/i18n"
	defaultLangHeader         = "Accept-Language"
	defaultLangResponseHeader = "Content-Language"
	defaultBundleFileFormat   = "json"
	defaultLang               = "en"
)

type Config struct {
	HeaderKey   string
	SourcePath  string
	Format      string
	DefaultLang string
}

// LocalizationMiddleware is a middleware that load message file and injects the localization bundle into the request context.
// Refer https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation
//
//	Precondition:
//	- Prepare language bundle files in the source path. (e.g. resources/i18n/en.json)
//	- The language header key is "Accept-Language" by default. (e.g. Accept-Language: en)
//	- The default language is "en" by default.
//	- The default bundle file format is "json" by default.
func LocalizationMiddleware(ctx context.Context, cfg Config) lit.HandlerFunc {
	cfg = prepareLocalizeConfig(cfg)
	bundle := i18n.Init(ctx, i18n.BundleConfig{
		DefaultLang:      cfg.DefaultLang,
		SourcePath:       cfg.SourcePath,
		BundleFileFormat: cfg.Format,
	})

	return func(c lit.Context) {
		req := c.Request()
		reqCtx := req.Context()

		// Load the message file only once for a given language.
		lang := req.Header.Get(cfg.HeaderKey)
		if lang != "" {
			lc := bundle.GetLocalize(lang)
			c.SetRequestContext(i18n.SetInContext(reqCtx, lc))
		}

		// Continue handle request
		c.Next()

		// Add localization header to response
		if lang == "" {
			lang = defaultLang
		}
		c.Writer().Header().Add(defaultLangResponseHeader, lang)
	}
}

func prepareLocalizeConfig(cfg Config) Config {
	if cfg.SourcePath == "" {
		cfg.SourcePath = defaultSourcePath
	}

	if cfg.Format == "" {
		cfg.Format = defaultBundleFileFormat
	}

	if cfg.HeaderKey == "" {
		cfg.HeaderKey = defaultLangHeader
	}

	if cfg.DefaultLang == "" {
		cfg.DefaultLang = defaultLang
	}

	return cfg
}
