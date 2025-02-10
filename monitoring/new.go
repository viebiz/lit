package monitoring

import (
	"io"
	"os"

	"go.uber.org/zap"
)

func NewLogger(opts ...Option) *Logger {
	return NewLoggerWithWriter(os.Stdout, opts...)
}

func NewLoggerWithWriter(w io.Writer, opts ...Option) *Logger {
	l := &Logger{
		zapLogger: zap.New(newZapCore(w)),
	}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func NewNoopLogger() *Logger {
	return &Logger{
		zapLogger: zap.NewNop(),
	}
}
