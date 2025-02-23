package monitoring

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/viebiz/lit/monitoring/tracing"
)

// Config holds Monitor configuration
type Config struct {
	ServerName      string
	Environment     string
	Version         string
	Writer          io.Writer // Support write log to buffer for testing
	SentryDSN       string    // To capture error, skip init Sentry if it's not provided
	OtelExporterURL string    // To support OpenTelemetry
	ExtraTags       map[string]string
}

// New creates a new Monitor instance
func New(cfg Config) (*Monitor, error) {
	// Setup logger
	var w io.Writer = os.Stdout
	if cfg.Writer != nil {
		w = cfg.Writer
	}

	m := &Monitor{
		logger:  zap.New(newZapCore(w)),
		logTags: map[string]string{},
	}

	if cfg.ExtraTags == nil {
		cfg.ExtraTags = make(map[string]string)
	}
	cfg.ExtraTags["server.name"] = cfg.ServerName
	cfg.ExtraTags["environment"] = cfg.Environment
	cfg.ExtraTags["version"] = cfg.Version
	m = m.With(cfg.ExtraTags)

	// Setup sentry
	sentryClient, err := initSentry(sentryConfig{
		DSN:         cfg.SentryDSN,
		ServerName:  cfg.ServerName,
		Environment: cfg.Environment,
		Version:     cfg.Version,
	}, m.logger)
	if err != nil {
		return nil, err
	}
	m.sentryClient = sentryClient

	// Setup tracing service
	if err := tracing.Init(context.Background(), tracing.Config{ExporterURL: cfg.OtelExporterURL}); err != nil {
		// Can skip if Exporter URL not provided
		if !errors.Is(err, tracing.ErrMissingExporterURL) {
			return nil, err
		}

		m.logger.Info("OTelExporter URL not provided. Not using Distributed Tracing")
	}

	return m, nil
}
