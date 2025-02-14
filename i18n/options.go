package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Option func(bundle *i18n.Bundle)

// WithUnmarshalFunc register new unmarshal function to load the message file
func WithUnmarshalFunc(format string, unmarshalFunc func(data []byte, v interface{}) error) Option {
	return func(bundle *i18n.Bundle) {
		bundle.RegisterUnmarshalFunc(format, unmarshalFunc)
	}
}
