package env

import (
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath  = "."
	defaultFormat      = "env"
	defaultEnvFilename = "config.env"
	envVarPrefix       = "app"
)

// ReadAppConfig loads app config with default settings:
//
//	Config Path: "."
//	Config Type: "env"
//	Config File: ".env"
func ReadAppConfig[T AppConfig]() (T, error) {
	return ReadAppConfigWithOptions[T](defaultFormat, defaultEnvFilename)
}

// ReadAppConfigWithOptions loads app config with custom options.
func ReadAppConfigWithOptions[T AppConfig](format string, fileName string, extraConfigPath ...string) (T, error) {
	// Create viper to read config
	v := viper.NewWithOptions()
	v.SetConfigType(format)
	v.AddConfigPath(defaultConfigPath)

	// Adding more config locations
	for _, loc := range extraConfigPath {
		v.AddConfigPath(loc)
	}

	v.SetConfigName(fileName)

	// Read config from file
	if err := v.ReadInConfig(); err != nil {
		return *new(T), pkgerrors.Wrap(err, "read from file error")
	}

	// Merge with ENV VARs
	v.SetEnvPrefix(envVarPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal to config object
	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return *new(T), pkgerrors.Wrap(err, "unmarshal to object error")
	}

	return cfg, nil
}
