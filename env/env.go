package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetAndValidate[T any](name string, validators ...Validator) (T, error) {
	var val T
	if err := parseFromString(strings.TrimSpace(os.Getenv(name)), &val); err != nil {
		return *new(T), err
	}

	for _, valid := range validators {
		if err := valid(name, val); err != nil {
			return *new(T), err
		}
	}

	return val, nil
}

func parseFromString[T any](value string, dst *T) error {
	switch interface{}(*dst).(type) {
	case string:
		*dst = interface{}(value).(T)

		return nil
	case int:
		if value == "" {
			return nil
		}

		parsed, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		*dst = interface{}(parsed).(T)
		return nil
	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		*dst = interface{}(parsed).(T)
		return nil
	default:
		return fmt.Errorf("env variable must be a string, number, or boolean")
	}
}

// Validator validates environment variable
type Validator func(name string, value interface{}) error

func Required() Validator {
	return func(name string, value interface{}) error {
		if value == nil {
			return fmt.Errorf("%s is required", name)
		}

		switch value.(type) {
		case string:
			if value == "" {
				return fmt.Errorf("%s is required", name)
			}
		}

		return nil
	}
}

func ExpectedValues(expected ...interface{}) Validator {
	return func(name string, value interface{}) error {
		for _, e := range expected {
			if value == e {
				return nil
			}
		}

		return fmt.Errorf("%s is not expected value: %s", name, value)
	}
}

func IsPositive() Validator {
	return func(name string, value interface{}) error {
		number, ok := value.(int)
		if !ok {
			return fmt.Errorf("%s is not a number", name)
		}

		if number < 0 {
			return fmt.Errorf("%s is not a positive number: %s", name, value)
		}

		return nil
	}
}
