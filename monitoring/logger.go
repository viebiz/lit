package monitoring

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	zapLogger *zap.Logger
}

func (l *Logger) With(fields ...LogField) *Logger {
	zapFields := make([]zap.Field, len(fields))
	for idx, field := range fields {
		zapFields[idx] = toZapField(field)
	}

	return &Logger{
		zapLogger: l.zapLogger.With(zapFields...),
	}
}

// WithLazy adds log fields lazily, where the field values are not computed
// when added, but are calculated when the log is written to output
func (l *Logger) WithLazy(fields ...LogField) *Logger {
	zapFields := make([]zap.Field, len(fields))
	for idx, field := range fields {
		zapFields[idx] = toZapField(field)
	}

	return &Logger{
		zapLogger: l.zapLogger.WithLazy(zapFields...),
	}
}

func (l *Logger) Infof(msg string, args ...any) {
	l.zapLogger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Errorf(err error, msg string, args ...any) {
	l.zapLogger.Error(fmt.Sprintf(msg, args...), zap.Error(err))

	// TODO: Log to sentry
}
