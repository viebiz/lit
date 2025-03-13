package i18n

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestBundle_LoadMessageFile(t *testing.T) {
	testCases := map[string]struct {
		sourcePath string
		langKey    string
		expNil     bool
		ext        string
		expErr     error
	}{
		"successful load": {
			sourcePath: "testdata",
			langKey:    "en",
			ext:        "json",
		},
		"file not found": {
			sourcePath: "INVALIDPATH",
			langKey:    "en",
			ext:        "json",
			expErr:     errors.New("stat INVALIDPATH: no such file or directory"),
		},
		"bundleMessage is nil": {
			sourcePath: "testdata",
			langKey:    "en",
			expNil:     true,
			ext:        "json",
		},
		"unsupported file format": {
			sourcePath: "testdata",
			langKey:    "en",
			ext:        "yaml",
			expErr:     errors.New("open testdata/en.yaml: no such file or directory"),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			b := &bundleMessage{
				underlyingBundle: i18n.NewBundle(language.English),
				DefaultLang:      "en",
				localizeMap:      make(map[string]Localizable),
			}
			b.underlyingBundle.RegisterUnmarshalFunc(defaultBundleFileFormat, json.Unmarshal)

			if tc.expNil {
				b = nil
			}

			// When
			err := b.LoadMessageFile(tc.sourcePath, tc.langKey, tc.ext)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Nil(t, b.localizeMap[tc.langKey])
			} else {
				require.NoError(t, err)
				if !tc.expNil {
					require.NotNil(t, b.localizeMap[tc.langKey])
				}
			}
		})
	}
}

func TestBundle_GetLocalize(t *testing.T) {
	tcs := map[string]struct {
		sourcePath     string
		langKey        string
		expAlreadyLoad bool
		expNonNil      bool
	}{
		"first load language": {
			sourcePath: "testdata",
			langKey:    "en",
			expNonNil:  true,
		},
		"already loaded language": {
			sourcePath:     "testdata",
			langKey:        "en",
			expAlreadyLoad: true,
			expNonNil:      true,
		},
		"invalid language": {
			sourcePath: "testdata",
			langKey:    "zz",
			expNonNil:  false,
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// GIVEN
			underlying := i18n.NewBundle(language.English)
			underlying.RegisterUnmarshalFunc(defaultBundleFileFormat, json.Unmarshal)
			b := &bundleMessage{
				DefaultLang:      "en",
				SourcePath:       tc.sourcePath,
				BundleFileFormat: defaultBundleFileFormat,
				underlyingBundle: underlying,
				langOnce:         sync.Map{},
				localizeMap:      make(map[string]Localizable),
			}

			if tc.expAlreadyLoad {
				b.localizeMap[tc.langKey] = newLocalizer(tc.langKey, b.underlyingBundle)
			}

			// WHEN
			lc := b.GetLocalize(tc.langKey)

			// THEN
			if tc.expNonNil {
				require.NotNil(t, lc, "expected non-nil localizer for language key: %s", tc.langKey)
			} else {
				require.Nil(t, lc, "expected nil localizer for language key: %s", tc.langKey)
			}
		})
	}
}
