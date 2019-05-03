package logutils

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFileSizeDefault = 10
	logFileNumBackups  = 1
	logFileAge         = 30
)

// Configure configures the shared default logger.
func Configure(level string, format string) zap.Config {

	return ConfigureWithOptions(level, format, "", false, false)
}

// ConfigureWithName configures the shared default logger.
func ConfigureWithName(serviceName string, level string, format string) zap.Config {

	logger, config := NewLogger(serviceName, level, format, "", false, false)

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)

	return config
}

// ConfigureWithOptions configures the shared default logger with options such as file and timestamp formats.
func ConfigureWithOptions(level string, format string, file string, fileOnly bool, prettyTimestamp bool) zap.Config {

	logger, config := NewLogger("", level, format, file, fileOnly, prettyTimestamp)

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)

	return config
}

// NewLogger returns a new configured zap.Logger
func NewLogger(serviceName string, level string, format string, file string, fileOnly bool, prettyTimestamp bool) (*zap.Logger, zap.Config) {

	var config zap.Config

	var initialFields map[string]interface{}
	if serviceName != "" {
		initialFields = map[string]interface{}{
			"srv": serviceName,
		}
	}

	switch format {
	case "json":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.EncoderConfig.CallerKey = "c"
		config.EncoderConfig.LevelKey = "l"
		config.EncoderConfig.MessageKey = "m"
		config.EncoderConfig.NameKey = "n"
		config.EncoderConfig.TimeKey = "t"

		config.InitialFields = initialFields

	case "stackdriver":
		config = zap.NewProductionConfig()
		config.EncoderConfig.LevelKey = "severity"
		config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch l {
			case zapcore.DebugLevel:
				enc.AppendString("DEBUG")
			case zapcore.InfoLevel:
				enc.AppendString("INFO")
			case zapcore.WarnLevel:
				enc.AppendString("WARNING")
			case zapcore.ErrorLevel:
				enc.AppendString("ERROR")
			case zapcore.DPanicLevel:
				enc.AppendString("CRITICAL")
			case zapcore.PanicLevel:
				enc.AppendString("ALERT")
			case zapcore.FatalLevel:
				enc.AppendString("EMERGENCY")
			}
		}

		config.InitialFields = initialFields

	default:
		config = zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Handle log file output
	w, err := handleOutputFile(&config, file, fileOnly)
	if err != nil {
		panic(err)
	}

	// Pretty timestamp
	if prettyTimestamp {
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
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

	logger := initLogger(w, config)

	if err != nil {
		panic(err)
	}

	return logger, config
}

// initLogger constructs the logger from the options
func initLogger(w zapcore.WriteSyncer, conf zap.Config) *zap.Logger {
	var enc zapcore.Encoder
	var coreFile zapcore.Core

	switch conf.Encoding {
	case "json":
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	case "console":
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	default:
		panic("unknown encoding")
	}

	console := zapcore.Lock(os.Stdout)
	coreConsole := zapcore.NewCore(enc, console, conf.Level)

	if w != nil {
		coreFile = zapcore.NewCore(enc, w, conf.Level)
	}

	core := zapcore.NewTee(
		coreFile,
		coreConsole,
	)

	logger := zap.New(core)
	return logger
}

// handleOutputFile handles options in log configs to redirect to file
func handleOutputFile(config *zap.Config, file string, fileOnly bool) (zapcore.WriteSyncer, error) {

	var w zapcore.WriteSyncer

	if file == "" {
		return nil, nil
	}
	dir := filepath.Dir(file)
	if dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if file != "" {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    logFileSizeDefault,
			MaxBackups: logFileNumBackups,
			MaxAge:     logFileAge,
		})
	}

	if fileOnly {
		config.OutputPaths = []string{file}
	} else {
		config.OutputPaths = append(config.OutputPaths, file)
	}

	return w, nil
}
