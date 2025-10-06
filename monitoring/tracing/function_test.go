package tracing

import (
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Test_validateConfig(t *testing.T) {
	type arg struct {
		givenCfg Config
		expErr   error
	}
	tcs := map[string]arg{
		"valid config - grpc": {
			givenCfg: Config{
				ExporterURL:   "localhost:4317",
				TransportType: TransportGRPC,
			},
			expErr: nil,
		},
		"valid config - http": {
			givenCfg: Config{
				ExporterURL:   "localhost:4318",
				TransportType: TransportHTTP,
			},
			expErr: nil,
		},
		"missing exporter url": {
			givenCfg: Config{
				TransportType: TransportGRPC,
			},
			expErr: ErrMissingExporterURL,
		},
		"empty exporter url": {
			givenCfg: Config{
				ExporterURL:   "",
				TransportType: TransportGRPC,
			},
			expErr: ErrMissingExporterURL,
		},
		"invalid transport type": {
			givenCfg: Config{
				ExporterURL:   "localhost:4317",
				TransportType: "invalid",
			},
			expErr: ErrInvalidTransportType,
		},
		"empty transport type": {
			givenCfg: Config{
				ExporterURL:   "localhost:4317",
				TransportType: "",
			},
			expErr: ErrInvalidTransportType,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			// When
			err := validateConfig(tc.givenCfg)

			// Then
			if tc.expErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_buildResource(t *testing.T) {
	type arg struct {
		givenCfg          Config
		expServiceName    string
		expServiceVersion string
		expEnvironment    string
	}
	tcs := map[string]arg{
		"all fields populated": {
			givenCfg: Config{
				ServerName:  "my-service",
				Environment: "production",
				Version:     "1.2.3",
			},
			expServiceName:    "my-service",
			expServiceVersion: "1.2.3",
			expEnvironment:    "production",
		},
		"minimal config": {
			givenCfg: Config{
				ServerName: "test-service",
			},
			expServiceName:    "test-service",
			expServiceVersion: "",
			expEnvironment:    "",
		},
		"empty config": {
			givenCfg:          Config{},
			expServiceName:    "",
			expServiceVersion: "",
			expEnvironment:    "",
		},
		"production environment": {
			givenCfg: Config{
				ServerName:  "api-gateway",
				Environment: "production",
				Version:     "2.0.0",
			},
			expServiceName:    "api-gateway",
			expServiceVersion: "2.0.0",
			expEnvironment:    "production",
		},
		"staging environment": {
			givenCfg: Config{
				ServerName:  "payment-service",
				Environment: "staging",
				Version:     "1.5.0-rc1",
			},
			expServiceName:    "payment-service",
			expServiceVersion: "1.5.0-rc1",
			expEnvironment:    "staging",
		},
		"development environment": {
			givenCfg: Config{
				ServerName:  "user-service",
				Environment: "development",
				Version:     "0.1.0-dev",
			},
			expServiceName:    "user-service",
			expServiceVersion: "0.1.0-dev",
			expEnvironment:    "development",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			// When
			resource := buildResource(tc.givenCfg)

			// Then
			require.NotNil(t, resource)

			// Extract attributes
			attrs := resource.Attributes()

			// Verify service name
			serviceName := ""
			for _, attr := range attrs {
				if attr.Key == semconv.ServiceNameKey {
					serviceName = attr.Value.AsString()
					break
				}
			}
			require.Equal(t, tc.expServiceName, serviceName)

			// Verify service version
			serviceVersion := ""
			for _, attr := range attrs {
				if attr.Key == semconv.ServiceVersionKey {
					serviceVersion = attr.Value.AsString()
					break
				}
			}
			require.Equal(t, tc.expServiceVersion, serviceVersion)

			// Verify environment
			environment := ""
			for _, attr := range attrs {
				if attr.Key == semconv.DeploymentEnvironmentKey {
					environment = attr.Value.AsString()
					break
				}
			}
			require.Equal(t, tc.expEnvironment, environment)

			// Verify telemetry SDK attributes are present
			foundTelemetrySDKName := false
			foundTelemetrySDKLanguage := false
			foundTelemetrySDKVersion := false

			for _, attr := range attrs {
				switch attr.Key {
				case semconv.TelemetrySDKNameKey:
					require.Equal(t, "opentelemetry", attr.Value.AsString())
					foundTelemetrySDKName = true
				case semconv.TelemetrySDKLanguageKey:
					require.Equal(t, "go", attr.Value.AsString())
					foundTelemetrySDKLanguage = true
				case semconv.TelemetrySDKVersionKey:
					// Just verify it's not empty
					require.NotEmpty(t, attr.Value.AsString())
					foundTelemetrySDKVersion = true
				}
			}

			require.True(t, foundTelemetrySDKName, "telemetry.sdk.name attribute should be present")
			require.True(t, foundTelemetrySDKLanguage, "telemetry.sdk.language attribute should be present")
			require.True(t, foundTelemetrySDKVersion, "telemetry.sdk.version attribute should be present")

			// Verify schema URL
			require.Equal(t, semconv.SchemaURL, resource.SchemaURL())
		})
	}
}

func Test_buildResource_Attributes(t *testing.T) {
	// Given
	cfg := Config{
		ServerName:  "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	// When
	resource := buildResource(cfg)

	// Then
	attrs := resource.Attributes()

	// Verify minimum expected attributes are present
	attrKeys := make(map[string]bool)
	for _, attr := range attrs {
		attrKeys[string(attr.Key)] = true
	}

	// Required attributes
	require.True(t, attrKeys["service.name"], "service.name should be present")
	require.True(t, attrKeys["service.version"], "service.version should be present")
	require.True(t, attrKeys["deployment.environment"], "deployment.environment should be present")
	require.True(t, attrKeys["telemetry.sdk.name"], "telemetry.sdk.name should be present")
	require.True(t, attrKeys["telemetry.sdk.language"], "telemetry.sdk.language should be present")
	require.True(t, attrKeys["telemetry.sdk.version"], "telemetry.sdk.version should be present")
}

func Test_buildResource_MultipleCallsIndependent(t *testing.T) {
	// Given
	cfg1 := Config{
		ServerName:  "service-1",
		Environment: "env-1",
		Version:     "1.0.0",
	}
	cfg2 := Config{
		ServerName:  "service-2",
		Environment: "env-2",
		Version:     "2.0.0",
	}

	// When
	resource1 := buildResource(cfg1)
	resource2 := buildResource(cfg2)

	// Then - Verify resources are independent
	attrs1 := resource1.Attributes()
	attrs2 := resource2.Attributes()

	serviceName1 := ""
	serviceName2 := ""

	for _, attr := range attrs1 {
		if attr.Key == semconv.ServiceNameKey {
			serviceName1 = attr.Value.AsString()
		}
	}

	for _, attr := range attrs2 {
		if attr.Key == semconv.ServiceNameKey {
			serviceName2 = attr.Value.AsString()
		}
	}

	require.Equal(t, "service-1", serviceName1)
	require.Equal(t, "service-2", serviceName2)
	require.NotEqual(t, serviceName1, serviceName2)
}
