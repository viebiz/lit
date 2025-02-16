package env

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAppConfig(t *testing.T) {
	type pgConfig struct {
		URL string `mapstructure:"URL"`
	}
	type webConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	type config struct {
		AppName  string    `mapstructure:"APP_NAME"`
		Lang     []string  `mapstructure:"LANG"`
		PGConfig pgConfig  `mapstructure:"DB"`
		Web      webConfig `mapstructure:"WEB"`
	}

	tcs := map[string]struct {
		configFile     string
		configFormat   string
		configLocation string
		extraEnv       map[string]string
		expResult      config
		expErr         error
	}{
		"success - file .env": {
			configFile:     "test.env",
			configFormat:   "env",
			configLocation: "testdata",
			expResult: config{
				AppName:  "lightning",
				Lang:     []string{"en", "vi"},
				PGConfig: pgConfig{URL: "postgres:thisisurl"},
				Web:      webConfig{Host: "0.0.0.0", Port: 8080},
			},
		},
		"success - file .yaml": {
			configFile:     "test.yaml",
			configFormat:   "yaml",
			configLocation: "testdata",
			expResult: config{
				AppName:  "lightning",
				Lang:     []string{"en", "vi"},
				PGConfig: pgConfig{URL: "postgres:thisisurl"},
				Web:      webConfig{Host: "0.0.0.0", Port: 8080},
			},
		},
		"success - merge env file .env": {
			configFile:     "test.env",
			configFormat:   "env",
			configLocation: "testdata",
			extraEnv: map[string]string{
				"APP_WEB_HOST": "192.168.0.1",
				"APP_LANG":     "vi,fr",
			},
			expResult: config{
				AppName:  "lightning",
				Lang:     []string{"vi", "fr"},
				PGConfig: pgConfig{URL: "postgres:thisisurl"},
				Web:      webConfig{Host: "192.168.0.1", Port: 8080},
			},
		},
		"success - merge env file .yaml": {
			configFile:     "test.yaml",
			configFormat:   "yaml",
			configLocation: "testdata",
			extraEnv: map[string]string{
				"APP_WEB_HOST": "192.168.0.1",
				"APP_LANG":     "vi,fr",
			},
			expResult: config{
				AppName:  "lightning",
				Lang:     []string{"vi", "fr"},
				PGConfig: pgConfig{URL: "postgres:thisisurl"},
				Web:      webConfig{Host: "192.168.0.1", Port: 8080},
			},
		},
		"error - config file not found": {
			configFile:     "nonexistent.env",
			configFormat:   "env",
			configLocation: "testdata",
			expErr:         errors.New("read from file error"),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			//t.Parallel() // Because override ENV VAR, so cannot run test parallel

			// Given
			for k, v := range tc.extraEnv {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tc.extraEnv {
					os.Unsetenv(k)
				}
			}()

			// When
			cfg, err := ReadAppConfigWithOptions[config](tc.configFormat, tc.configFile, tc.configLocation)

			// Then
			if tc.expErr != nil {
				require.ErrorContains(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, cfg)
			}
		})
	}
}

func TestReadAppConfig_UnmarshalError(t *testing.T) {
	// Given
	type config struct {
		AppName chan int `mapstructure:"APP_NAME"` // A channel can't be unmarshaled, forcing an error
	}

	// When
	cfg, err := ReadAppConfig[config]()

	// Then
	require.Zero(t, cfg)
	require.ErrorContains(t, err, "unmarshal to object error")
}
