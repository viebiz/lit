package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetMode_IsValid(t *testing.T) {
	type args struct {
		givenMode setMode
		expValid  bool
	}
	tcs := map[string]args{
		"valid NX": {
			givenMode: setModeNX,
			expValid:  true,
		},
		"valid XX": {
			givenMode: setModeXX,
			expValid:  true,
		},
		"invalid none": {
			givenMode: setModeNone,
			expValid:  false,
		},
		"invalid empty": {
			givenMode: "",
			expValid:  false,
		},
		"invalid random": {
			givenMode: "INVALID",
			expValid:  false,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			result := tc.givenMode.IsValid()

			// Then
			require.Equal(t, tc.expValid, result)
		})
	}
}

func TestSetMode_String(t *testing.T) {
	type args struct {
		givenMode setMode
		expString string
	}
	tcs := map[string]args{
		"NX": {
			givenMode: setModeNX,
			expString: "NX",
		},
		"XX": {
			givenMode: setModeXX,
			expString: "XX",
		},
		"none": {
			givenMode: setModeNone,
			expString: "",
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			result := tc.givenMode.String()

			// Then
			require.Equal(t, tc.expString, result)
		})
	}
}
