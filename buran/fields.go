package buran

import (
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap/zapcore"
)

func zapFieldsToOpentracing(zapFields ...zapcore.Field) []opentracinglog.Field {
	opentracingFields := make([]opentracinglog.Field, len(zapFields))

	for i, zapField := range zapFields {
		opentracingFields[i] = zapFieldToOpentracing(zapField)
	}

	return opentracingFields
}

func zapFieldToOpentracing(zapField zapcore.Field) opentracinglog.Field {
	switch zapField.Type {

	case zapcore.BoolType:
		return opentracinglog.Bool(zapField.Key, zapField.Interface.(bool))
	case zapcore.Float32Type:
		return opentracinglog.Float32(zapField.Key, zapField.Interface.(float32))
	case zapcore.Float64Type:
		return opentracinglog.Float64(zapField.Key, zapField.Interface.(float64))
	case zapcore.Int64Type:
		return opentracinglog.Int64(zapField.Key, zapField.Interface.(int64))
	case zapcore.Int32Type:
		return opentracinglog.Int32(zapField.Key, zapField.Interface.(int32))
	case zapcore.StringType:
		return opentracinglog.String(zapField.Key, zapField.Interface.(string))
	case zapcore.Uint64Type:
		return opentracinglog.Uint64(zapField.Key, zapField.Interface.(uint64))
	case zapcore.Uint32Type:
		return opentracinglog.Uint32(zapField.Key, zapField.Interface.(uint32))
	case zapcore.ErrorType:
		return opentracinglog.Error(zapField.Interface.(error))
	default:
		return opentracinglog.Object(zapField.Key, zapField.Interface)
	}
}
