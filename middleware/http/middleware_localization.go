package http

import (
	"sync"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/i18n"
	"github.com/viebiz/lit/monitoring"
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
	AcceptLang  []string
}

// LocalizationMiddleware is a middleware that load message file and injects the localization bundle into the request context.
// Refer https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation
//
//	Precondition:
//	- Prepare language bundle files in the source path. (e.g. resources/i18n/en.json)
//	- The language header key is "Accept-Language" by default. (e.g. Accept-Language: en)
//	- The default language is "en" by default.
//	- The default bundle file format is "json" by default.
func LocalizationMiddleware(cfg Config) lit.HandlerFunc {
	cfg = prepareLocalizeConfig(cfg)
	bundle := i18n.Init(i18n.BundleConfig{})

	// Initialize sync.Map with a sync.Once per accepted language.
	var langOnce sync.Map
	for _, lang := range cfg.AcceptLang {
		langOnce.Store(lang, &sync.Once{})
	}

	return func(c lit.Context) {
		req := c.Request()
		ctx := req.Context()
		lang := req.Header.Get(cfg.HeaderKey)

		// Load the message file only once for a given language.
		if lang != "" {
			if onceInterface, ok := langOnce.Load(lang); ok {
				once := onceInterface.(*sync.Once)
				once.Do(func() {
					if err := bundle.LoadMessageFile(cfg.SourcePath, lang, cfg.Format); err != nil {
						monitoring.FromContext(ctx).Errorf(err, "Failed to load message file")
					}
				})
			}

		}

		// Inject localization bundle to request context
		c.SetRequestContext(i18n.SetInContext(ctx, bundle))

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
