package i18n

import (
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestTranslator_Translate(t *testing.T) {
	tcs := map[string]struct {
		givenMsgID    string
		givenMsgParam map[string]any
		expResult     string
		expErr        error
	}{
		"success - message": {
			givenMsgID: "success",
			expResult:  "Success",
		},
		"success - message with params": {
			givenMsgID: "helloPerson",
			givenMsgParam: map[string]any{
				"Name": "The Knight",
			},
			expResult: "Hello The Knight",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			tr, err := NewTranslator([]string{"en", "vi"}, "testdata")
			require.NoError(t, err)

			// When
			out, err := tr.Translate(tc.givenMsgID, tc.givenMsgParam)

			// Then
			require.NoError(t, err)
			require.Equal(t, tc.expResult, out)
		})
	}
}

func TestTranslator_TranslateWithLang(t *testing.T) {
	tcs := map[string]struct {
		lang          string
		givenMsgID    string
		givenMsgParam map[string]any
		expResult     string
		expErr        error
	}{
		"success - message in english": {
			lang:       "en",
			givenMsgID: "success",
			expResult:  "Success",
		},
		"success - message with params in english": {
			lang:       "en",
			givenMsgID: "helloPerson",
			givenMsgParam: map[string]any{
				"Name": "The Knight",
			},
			expResult: "Hello The Knight",
		},
		"success - message in vietnamese": {
			lang:       "vi",
			givenMsgID: "success",
			expResult:  "Thành công",
		},
		"success - message with params in vietnamese": {
			lang:       "vi",
			givenMsgID: "helloPerson",
			givenMsgParam: map[string]any{
				"Name": "The Knight",
			},
			expResult: "Chào The Knight",
		},
		"error - lang not supported": {
			lang:       "zh",
			givenMsgID: "success",
			expErr:     ErrGivenLangNotSupported,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			tr, err := NewTranslator([]string{"en", "vi"}, "testdata")
			require.NoError(t, err)

			// When
			out, err := tr.TranslateWithLang(tc.lang, tc.givenMsgID, tc.givenMsgParam)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, pkgerrors.Unwrap(err), tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, out)
			}
		})
	}
}
