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

func WithFieldFromMap(fields map[string]string) Option {
	return func(l *Logger) {
		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.String(k, v))
		}

		l.zapLogger = l.zapLogger.With(zapFields...)
	}
}
