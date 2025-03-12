package i18n

import (
	"errors"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestMessageLocalize(t *testing.T) {
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
			expResult: "",
			expErr:    ErrBundleNotInitialized,
		},
		"localize error": {
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
			simulateLocalizeErr: true,
			messageID:           "greeting",
			params:              nil,
			expErr:              errors.New("message id \"greeting\" does not match default message id \"INVALID\""),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given
			if tc.simulateLocalizeErr {
				defer func() { prepareMessageConfigWrapper = prepareMessageConfig }()

				prepareMessageConfigWrapper = func(messageID string, params map[string]interface{}) *i18n.LocalizeConfig {
					return &i18n.LocalizeConfig{MessageID: messageID, DefaultMessage: &i18n.Message{ID: "INVALID"}}
				}
			}

			// When
			ml := messageLocalize{
				localizer: tc.localizer,
			}

			// Then
			res, err := ml.Localize(tc.messageID, tc.params)
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, res)
			}
		})
	}
}
