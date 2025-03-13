package i18n

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestInit(t *testing.T) {
	tcs := map[string]struct {
		cfg                 BundleConfig
		expectedDefaultLang string
	}{
		"Default config": {
			cfg:                 BundleConfig{},
			expectedDefaultLang: defaultLang,
		},
		"Custom DefaultLanguage and AcceptLanguage": {
			cfg: BundleConfig{
				DefaultLang: language.Spanish.String(),
			},
			expectedDefaultLang: language.Spanish.String(),
		},
		"Custom BundleFileFormat and UnmarshalFunc": {
			cfg: BundleConfig{
				ExtraBundleFileSupport: map[string]UnmarshalFunc{
					"yaml": func(data []byte, v interface{}) error {
						// Dummy unmarshal: use json.Unmarshal here for testing
						return json.Unmarshal(data, v)
					},
				},
			},
			expectedDefaultLang: defaultLang,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			bundleInterface := Init(context.Background(), tc.cfg)

			// Then
			require.NotNil(t, bundleInterface)

			b, ok := bundleInterface.(*bundleMessage)
			require.True(t, ok, "returned bundleMessage is not of expected type")
			require.Equal(t, tc.expectedDefaultLang, b.DefaultLang, "default language not set as expected")
			require.NotNil(t, b.underlyingBundle)
		})
	}
}
