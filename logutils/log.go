package logutils

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configure configures the shared default logger.
func Configure(level string, format string) zap.Config {

	var config zap.Config

	switch format {
	case "json":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	default:
		config = zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set the logger
	switch level {
	case "trace", "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)

	return config
}