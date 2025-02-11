package monitoring

import (
	"go.uber.org/zap"
)

type Option func(*Logger)

func WithFields(fields ...LogField) Option {
	return func(l *Logger) {
		zapFields := make([]zap.Field, len(fields))
		for idx, f := range fields {
			zapFields[idx] = toZapField(f)
		}

		l.zapLogger = l.zapLogger.With(zapFields...)
	}
}
