package monitoring

import (
	"context"
	"errors"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapCore(w io.Writer) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // to show {"level": "info"} or to show {"level": "INFO"}
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // to format time as {"timestamp":"2021-06-21T09:25:51.230+08:00"}
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// flush flushes any buffered log entries.
func flushZap(z *zap.Logger, maxWait time.Duration) error {
	errChan := make(chan error)
	go func() {
		errChan <- z.Sync()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()
	select {
	case <-errChan:
		// NOTE: We ignore any errors here because Sync is known to fail with EINVAL
		// When logging to Stdout on certain OS's.
		//
		// Uber made the same change within the core of the lg implementation.
		// See: https://github.com/uber-go/zap/issues/328
		// See: https://github.com/influxdata/influxdb/pull/20448
		return nil
	case <-ctx.Done():
		return errors.New("timed out")
	}
}
