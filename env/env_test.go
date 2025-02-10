package env

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getAndValidEnv(t *testing.T) {
	type args struct {
		Name    string
		Age     int
		Married bool
	}
	tcs := map[string]struct {
		givenInput map[string]string
		config     map[string][]Validator
		expResult  args
		expErr     error
	}{
		"success": {
			givenInput: map[string]string{
				"NAME":    "vae",
				"AGE":     "25",
				"MARRIED": "false",
			},
			config: map[string][]Validator{
				"NAME":    {Required()},
				"AGE":     {Required(), IsPositive()},
				"MARRIED": {Required()},
			},
			expResult: args{
				Name:    "vae",
				Age:     25,
				Married: false,
			},
		},
		"success without validation": {
			givenInput: map[string]string{
				"NAME":    "vae",
				"AGE":     "25",
				"MARRIED": "false",
			},
			expResult: args{
				Name:    "vae",
				Age:     25,
				Married: false,
			},
		},
		"failure - fail validation": {
			config: map[string][]Validator{
				"NAME": {Required()},
			},
			expErr: errors.New("NAME is required"),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			for k, v := range tc.givenInput {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tc.givenInput {
					os.Unsetenv(k)
				}
			}()

			// Prepare func to get env
			readFunc := func(values map[string]string) (args, error) {
				name, err := GetAndValidate[string]("NAME", tc.config["NAME"]...)
				if err != nil {
					return args{}, err
				}

				age, err := GetAndValidate[int]("AGE", tc.config["AGE"]...)
				if err != nil {
					return args{}, err
				}

				married, err := GetAndValidate[bool]("MARRIED", tc.config["MARRIED"]...)
				if err != nil {
					return args{}, err
				}

				return args{
					Name:    name,
					Age:     age,
					Married: married,
				}, nil
			}

			// When
			out, err := readFunc(tc.givenInput)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, out)
			}
		})
	}
}
