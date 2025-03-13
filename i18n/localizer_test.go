package i18n

import (
	"errors"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestLocalizer_TryLocalize(t *testing.T) {
	tcs := map[string]struct {
		localizer           *i18n.Localizer
		messageID           string
		params              map[string]interface{}
		simulateLocalizeErr bool
		expResult           string
		expErr              error
	}{
		"success without params": {
			localizer: func() *i18n.Localizer {
				bundle := i18n.NewBundle(language.English)
				// Add a simple message.
				err := bundle.AddMessages(language.English, &i18n.Message{
					ID:    "greeting",
					Other: "Hello",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(bundle, "en")
			}(),
			messageID: "greeting",
			expResult: "Hello",
		},
		"success with params": {
			localizer: func() *i18n.Localizer {
				bundle := i18n.NewBundle(language.English)
				// Add a message with a template.
				err := bundle.AddMessages(language.English, &i18n.Message{
					ID:    "welcome",
					Other: "Hello, {{.Name}}",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(bundle, "en")
			}(),
			messageID: "welcome",
			params:    map[string]interface{}{"Name": "John"},
			expResult: "Hello, John",
		},
		"nil localizer": {
			localizer: nil,
			messageID: "anything",
			params:    nil,
			expResult: "anything",
		},
		"localize error": {
			localizer: func() *i18n.Localizer {
				b := i18n.NewBundle(language.English)
				// Add a simple message.
				err := b.AddMessages(language.English, &i18n.Message{
					ID:    "greeting",
					Other: "Hello",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(b, "en")
			}(),
			simulateLocalizeErr: true,
			messageID:           "greeting",
			params:              nil,
			expErr:              errors.New("simulated localize error"),
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			l := localizer{
				underlyingLocalizer: tc.localizer,
				getLocalizedMessage: getLocalizedMessage,
			}

			if tc.simulateLocalizeErr {
				l.getLocalizedMessage = func(localizer *i18n.Localizer, cfg *i18n.LocalizeConfig) (string, error) {
					return "", errors.New("simulated localize error")
				}
			}

			// When
			res, err := l.TryLocalize(tc.messageID, tc.params)

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

func TestLocalizer_Localize(t *testing.T) {
	tcs := map[string]struct {
		localizer           *i18n.Localizer
		messageID           string
		params              map[string]interface{}
		simulateLocalizeErr bool
		expResult           string
		expErr              error
	}{
		"success without params": {
			localizer: func() *i18n.Localizer {
				bundle := i18n.NewBundle(language.English)
				// Add a simple message.
				err := bundle.AddMessages(language.English, &i18n.Message{
					ID:    "greeting",
					Other: "Hello",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(bundle, "en")
			}(),
			messageID: "greeting",
			expResult: "Hello",
		},
		"success with params": {
			localizer: func() *i18n.Localizer {
				bundle := i18n.NewBundle(language.English)
				// Add a message with a template.
				err := bundle.AddMessages(language.English, &i18n.Message{
					ID:    "welcome",
					Other: "Hello, {{.Name}}",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(bundle, "en")
			}(),
			messageID: "welcome",
			params:    map[string]interface{}{"Name": "John"},
			expResult: "Hello, John",
		},
		"nil localizer": {
			localizer: nil,
			messageID: "anything",
			params:    nil,
			expResult: "anything",
		},
		"localize error": {
			localizer: func() *i18n.Localizer {
				b := i18n.NewBundle(language.English)
				// Add a simple message.
				err := b.AddMessages(language.English, &i18n.Message{
					ID:    "greeting",
					Other: "Hello",
				})
				require.NoError(t, err)
				return i18n.NewLocalizer(b, "en")
			}(),
			simulateLocalizeErr: true,
			messageID:           "greeting",
			params:              nil,
			expErr:              errors.New("simulated localize error"),
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			l := localizer{
				underlyingLocalizer: tc.localizer,
				getLocalizedMessage: getLocalizedMessage,
			}

			if tc.simulateLocalizeErr {
				l.getLocalizedMessage = func(localizer *i18n.Localizer, cfg *i18n.LocalizeConfig) (string, error) {
					return "", errors.New("simulated localize error")
				}
			}

			// When
			res := l.Localize(tc.messageID, tc.params)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.messageID, res)
			} else {
				require.Equal(t, tc.expResult, res)
			}
		})
	}
}
