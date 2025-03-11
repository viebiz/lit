package i18n

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestInit(t *testing.T) {
	testCases := map[string]struct {
		cfg                 BundleConfig
		expectedDefaultLang string
	}{
		"Default config": {
			cfg:                 BundleConfig{},
			expectedDefaultLang: defaultLangTag,
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
			expectedDefaultLang: defaultLangTag,
		},
	}

	for name, tc := range testCases {
		tc := tc // capture range variable
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			bundleInterface := Init(tc.cfg)

			// Then
			require.NotNil(t, bundleInterface)

			b, ok := bundleInterface.(*bundle)
			require.True(t, ok, "returned bundle is not of expected type")
			require.Equal(t, tc.expectedDefaultLang, b.DefaultLang, "default language not set as expected")
			require.NotNil(t, b.i18nBundle)
		})
	}
}
