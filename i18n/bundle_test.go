package i18n

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestBundle_LoadMessageFile_Table(t *testing.T) {
	testCases := map[string]struct {
		sourcePath string
		langKey    string
		ext        string
		expErr     error
	}{
		"successful load": {
			// file is already prepared in the testdata folder
			sourcePath: "testdata",
			langKey:    "en",
			ext:        "json",
		},
		"file not found": {
			// invalid source path to trigger error since file will not be found
			sourcePath: "INVALIDPATH",
			langKey:    "en",
			ext:        "json",
			expErr:     errors.New("stat INVALIDPATH: no such file or directory"),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			b := &bundle{
				i18nBundle:  i18n.NewBundle(language.English),
				DefaultLang: "en",
				LocalizeMap: make(map[string]MessageLocalize),
			}
			b.i18nBundle.RegisterUnmarshalFunc(defaultBundleFileFormat, json.Unmarshal)

			// When
			err := b.LoadMessageFile(tc.sourcePath, tc.langKey, tc.ext)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())

				_, exists := b.LocalizeMap[tc.langKey]
				require.False(t, exists, "localize map should not contain the language key on error")
			} else {
				require.NoError(t, err)
				_, exists := b.LocalizeMap[tc.langKey]
				require.True(t, exists, "localize map should contain the language key")
			}
		})
	}
}

func TestLocaleManager_LocalizeWithLang(t *testing.T) {
	type mockInfo struct {
		useMock   bool
		returnMsg string
		returnErr error
	}
	tcs := map[string]struct {
		defaultLang string
		langKey     string
		messageID   string
		params      map[string]interface{}
		mock        mockInfo
		expResult   string
		expErr      error
	}{
		"unsupported language": {
			defaultLang: "en",
			langKey:     "fr",
			messageID:   "greeting",
			params:      nil,
			expResult:   "",
			expErr:      ErrGivenLangNotSupported,
		},
		"localizer returns error": {
			defaultLang: "en",
			langKey:     "fr",
			messageID:   "greeting",
			params:      nil,
			mock:        mockInfo{useMock: true, returnMsg: "", returnErr: errors.New("localization error")},
			expResult:   "",
			expErr:      errors.New("localization error"),
		},
		"successful localization": {
			defaultLang: "en",
			langKey:     "fr",
			messageID:   "greeting",
			params:      map[string]interface{}{"key": "value"},
			mock:        mockInfo{useMock: true, returnMsg: "Bonjour", returnErr: nil},
			expResult:   "Bonjour",
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			localizeBundle := make(map[string]MessageLocalize)
			if tc.mock.useMock {
				mockLocalizer := NewMockMessageLocalize(t)
				mockLocalizer.
					EXPECT().
					Localize(tc.messageID, tc.params).
					Return(tc.mock.returnMsg, tc.mock.returnErr)
				localizeBundle[tc.langKey] = mockLocalizer
			}
			lm := bundle{
				DefaultLang: tc.defaultLang,
				LocalizeMap: localizeBundle,
			}

			// When
			res, err := lm.LocalizeWithLang(tc.langKey, tc.messageID, tc.params)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, res)
			}
		})
	}
}

func TestLocaleManager_Localize(t *testing.T) {
	type mockInfo struct {
		useMock   bool
		returnMsg string
		returnErr error
	}

	tcs := map[string]struct {
		defaultLang string
		messageID   string
		params      map[string]interface{}
		mock        mockInfo
		expResult   string
		expErr      error
	}{
		"default language unsupported": {
			defaultLang: "fr", // No localizer provided for "fr"
			messageID:   "welcome",
			params:      nil,
			expResult:   "",
			expErr:      ErrGivenLangNotSupported,
		},
		"localizer returns error": {
			defaultLang: "en",
			messageID:   "welcome",
			params:      nil,
			mock:        mockInfo{useMock: true, returnMsg: "", returnErr: errors.New("default language error")},
			expResult:   "",
			expErr:      errors.New("default language error"),
		},
		"default localization success": {
			defaultLang: "en",
			messageID:   "welcome",
			params:      nil,
			mock:        mockInfo{useMock: true, returnMsg: "Hello", returnErr: nil},
			expResult:   "Hello",
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			localizeBundle := make(map[string]MessageLocalize)
			if tc.mock.useMock {
				mockLocalizer := NewMockMessageLocalize(t)
				mockLocalizer.
					EXPECT().
					Localize(tc.messageID, tc.params).
					Return(tc.mock.returnMsg, tc.mock.returnErr)
				localizeBundle[tc.defaultLang] = mockLocalizer
			}

			lm := bundle{
				DefaultLang: tc.defaultLang,
				LocalizeMap: localizeBundle,
			}

			// When
			res, err := lm.Localize(tc.messageID, tc.params)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, res)
			}
		})
	}
}

func TestLocaleManager_TryLocalize(t *testing.T) {
	type mockInfo struct {
		useMock   bool
		returnMsg string
		returnErr error
	}

	tcs := map[string]struct {
		defaultLang string
		messageID   string
		params      map[string]interface{}
		mock        mockInfo
		expResult   string
	}{
		"try localize success": {
			defaultLang: "en",
			messageID:   "salutation",
			params:      nil,
			mock:        mockInfo{useMock: true, returnMsg: "Hi", returnErr: nil},
			expResult:   "Hi",
		},
		"try localize fallback on error": {
			defaultLang: "en",
			messageID:   "salutation",
			params:      nil,
			mock:        mockInfo{useMock: true, returnMsg: "", returnErr: errors.New("translation missing")},
			expResult:   "salutation",
		},
		"try localize without localizer": {
			defaultLang: "en",
			messageID:   "salutation",
			params:      nil,
			expResult:   "salutation",
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			localizeBundle := make(map[string]MessageLocalize)
			if tc.mock.useMock {
				mockLocalizer := NewMockMessageLocalize(t)
				mockLocalizer.
					EXPECT().
					Localize(tc.messageID, tc.params).
					Return(tc.mock.returnMsg, tc.mock.returnErr)
				localizeBundle[tc.defaultLang] = mockLocalizer
			}

			lm := bundle{
				DefaultLang: tc.defaultLang,
				LocalizeMap: localizeBundle,
			}

			// When
			res := lm.TryLocalize(tc.messageID, tc.params)

			// Then
			require.Equal(t, tc.expResult, res)
		})
	}
}
