package monitoring

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogField zap.Field

func (l LogField) toAttribute() attribute.KeyValue {
	switch l.Type {
	case zapcore.StringType:
		return attribute.String(l.Key, l.String)
	case zapcore.Int64Type:
		return attribute.Int64(l.Key, l.Integer)
	case zapcore.ByteStringType:
		return attribute.String(l.Key, string(l.Interface.([]byte)))
	case zapcore.BoolType:
		return attribute.Bool(l.Key, l.Interface.(bool))
	default:
		if l.Interface != nil {
			return attribute.String(l.Key, fmt.Sprint(l.Interface))
		}

		if l.Integer != 0 {
			return attribute.Int64(l.Key, l.Integer)
		}

		if l.String != "" {
			return attribute.String(l.Key, l.String)
		}

		return attribute.KeyValue{}
	}
}

func Field[T any](name string, value T) LogField {
	switch v := interface{}(value).(type) {
	case string:
		return LogField(zap.String(name, v))
	case []string:
		return LogField(zap.Strings(name, v))
	case int:
		return LogField(zap.Int(name, v))
	case []byte:
		return LogField(zap.ByteString(name, v))
	case bool:
		return LogField(zap.Bool(name, v))
	default:
		return LogField(zap.Any(name, v))
	}
}

func toZapField(f LogField) zap.Field {
	return zap.Field(f)
}
