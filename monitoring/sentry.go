package monitoring

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

type sentryConfig struct {
	DSN         string
	ServerName  string
	Environment string
	Version     string
}

func initSentry(cfg sentryConfig, logger *zap.Logger) (*sentry.Client, error) {
	if cfg.DSN == "" {
		logger.Info("Sentry DSN not provided. Not using Sentry Error Reporting")
		return nil, nil
	}

	client, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:              cfg.DSN,
			AttachStacktrace: true,
			SampleRate:       1, // send all events
			// Integrations: func(integrations []sentry.Integration) []sentry.Integration { // may need to enable this if ever we go to Sentry cloud so as to not expose our source code.
			// 	var filteredIntegrations []sentry.Integration
			// 	for _, integration := range integrations {
			// 		if integration.Name() == "ContextifyFrames" {
			// 			continue
			// 		}
			// 		filteredIntegrations = append(filteredIntegrations, integration)
			// 	}
			// 	return filteredIntegrations
			// },
			ServerName:  cfg.ServerName,
			Release:     cfg.Version,
			Environment: cfg.Environment,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create sentry client err: %w", err)
	}

	logger.Info("Sentry Error Reporter initialized")

	return client, nil
}
