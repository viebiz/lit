package i18n

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTranslator(t *testing.T) {
	tcs := map[string]struct {
		givenSupportedLang []string
		givenBasePath      string
		expErr             error
	}{
		"success": {
			givenSupportedLang: []string{"en"},
			givenBasePath:      "testdata",
		},
		"error - message path not exists": {
			givenSupportedLang: []string{"en"},
			givenBasePath:      "NOTEXISTPATH",
			expErr:             errors.New("message path does not exist"),
		},
		"error - load message file": {
			givenSupportedLang: []string{"fr"},
			givenBasePath:      "testdata",
			expErr:             errors.New("open testdata/fr.json: no such file or directory"),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			rs, err := NewTranslator(tc.givenSupportedLang, tc.givenBasePath)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.NotZero(t, rs)
			}
		})
	}
}
