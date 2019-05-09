// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logutils

import (
	"errors"
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

	logger, config := newLogger(serviceName, level, format, "", false, false)

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)

	return config
}

// ConfigureWithOptions configures the shared default logger with options such as file and timestamp formats.
func ConfigureWithOptions(level string, format string, file string, fileOnly bool, prettyTimestamp bool) zap.Config {

	logger, config := newLogger("", level, format, file, fileOnly, prettyTimestamp)

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)

	return config
}

// newLogger returns a new configured zap.Logger
func newLogger(serviceName string, level string, format string, file string, fileOnly bool, prettyTimestamp bool) (*zap.Logger, zap.Config) {

	config := getConfig(serviceName, format)

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
	config.Level = levelToZapLevel(level)

	logger, err := config.Build()
	if w != nil {
		logger, err = config.Build(SetOutput(w, config))
	}
	if err != nil {
		panic(err)
	}

	return logger, config
}

// getConfig provides a zap configuration
func getConfig(serviceName, format string) zap.Config {

	var initialFields map[string]interface{}
	if serviceName != "" {
		initialFields = map[string]interface{}{
			"srv": serviceName,
		}
	}

	switch format {
	case "json":
		return getJSONConfig(initialFields)

	case "stackdriver":
		return getStackdriverConfig(initialFields)

	default:
		return getDefaultConfig()
	}
}

// getJSONConfig provides a JSON zap configuration
func getJSONConfig(initialFields map[string]interface{}) zap.Config {

	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.CallerKey = "c"
	config.EncoderConfig.LevelKey = "l"
	config.EncoderConfig.MessageKey = "m"
	config.EncoderConfig.NameKey = "n"
	config.EncoderConfig.TimeKey = "t"

	config.InitialFields = initialFields

	return config
}

// getStackdriverConfig provides a stackdriver zap configuration
func getStackdriverConfig(initialFields map[string]interface{}) zap.Config {

	config := zap.NewProductionConfig()
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

	return config
}

func getDefaultConfig() zap.Config {

	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return config
}

// levelToZapLevel provides a zapLevel given a level configuration
func levelToZapLevel(level string) zap.AtomicLevel {

	switch level {
	case "trace", "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

// getEncoder provides an encoder based on encoding configuration
func getEncoder(c zap.Config) (zapcore.Encoder, error) {

	switch c.Encoding {
	case "json":
		return zapcore.NewJSONEncoder(c.EncoderConfig), nil
	case "console":
		return zapcore.NewConsoleEncoder(c.EncoderConfig), nil
	default:
		return nil, errors.New("unknown encoding")
	}
}

// SetOutput returns the zap option with the new sync writer
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
func handleOutputFile(config *zap.Config, file string, fileOnly bool) (zapcore.WriteSyncer, error) {

	if file == "" {
		return nil, nil
	}
	dir := filepath.Dir(file)
	if dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	var w zapcore.WriteSyncer
	w = zapcore.AddSync(&lumberjack.Logger{
		Filename:   file,
		MaxSize:    logFileSizeDefault,
		MaxBackups: logFileNumBackups,
		MaxAge:     logFileAge,
	})

	if fileOnly {
		config.OutputPaths = []string{file}
		return w, nil
	}

	config.OutputPaths = append(config.OutputPaths, file)
	return w, nil
}
