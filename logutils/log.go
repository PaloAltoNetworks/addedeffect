package logutils

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configure configures the shared default logger.
func Configure(level string, format string) zap.Config {

	return ConfigureWithOptions(level, format, "", false, false)
}

// ConfigureWithOptions configures the shared default logger with options such as file and timestamp formats.
func ConfigureWithOptions(level string, format string, file string, fileOnly bool, prettyTimestamp bool) zap.Config {

	var config zap.Config

	switch format {
	case "json":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.EncoderConfig.CallerKey = "c"
		config.EncoderConfig.LevelKey = "l"
		config.EncoderConfig.MessageKey = "m"
		config.EncoderConfig.NameKey = "n"
		config.EncoderConfig.TimeKey = "t"

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
	default:
		config = zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Handle log file output
	if err := handleOutputFile(&config, file, fileOnly); err != nil {
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
	if fileOnly == true || file != "" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    1,
			MaxBackups: 3,
			MaxAge:     8,
		})
		logger, err := config.Build(SetOutput(w, config))
		if err != nil {
			panic(err)
		}

		zap.ReplaceGlobals(logger)
	} else {
		logger, err := config.Build()

		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(logger)

	}
	go handleElevationSignal(config)

	return config
}

func SetOutput(w zapcore.WriteSyncer, conf zap.Config) zap.Option {
	var enc zapcore.Encoder
	// Copy paste from zap.Config.buildEncoder.
	switch conf.Encoding {
	case "json":
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	case "console":
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	default:
		panic("unknown encoding")
	}
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(enc, w, conf.Level)
	})
}

// handleOutputFile handles options in log configs to redirect to file
func handleOutputFile(config *zap.Config, file string, fileOnly bool) error {

	if file == "" {
		return nil
	}

	dir := filepath.Dir(file)
	if dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	if fileOnly {
		config.OutputPaths = []string{file}
	} else {
		config.OutputPaths = append(config.OutputPaths, file)
	}

	return nil
}
