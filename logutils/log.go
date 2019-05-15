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
	logFileSizeDefault = 1
	logFileNumBackups  = 1
	logFileAge         = 30
)

// getConfig provides a zap configuration
func getConfig(level, format string) (config zap.Config) {

	switch format {
	case "json":
		config = getJSONConfig()

	case "stackdriver":
		config = getStackdriverConfig()

	default:
		config = getDefaultConfig()
	}

	// Set the logger
	config.Level = levelToZapLevel(level)
	return config
}

// getJSONConfig provides a JSON zap configuration
func getJSONConfig() zap.Config {

	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.CallerKey = "c"
	config.EncoderConfig.LevelKey = "l"
	config.EncoderConfig.MessageKey = "m"
	config.EncoderConfig.NameKey = "n"
	config.EncoderConfig.TimeKey = "t"

	return config
}

// getStackdriverConfig provides a stackdriver zap configuration
func getStackdriverConfig() zap.Config {

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

// initLogger constructs the logger from the options
func initLogger(conf *Config) (*zap.Logger, error) {

	enc, err := getEncoder(conf.z)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(enc, zapcore.Lock(os.Stderr), conf.z.Level)
	if conf.w != nil {
		if conf.fileOnly {
			core = zapcore.NewCore(enc, conf.w, conf.z.Level)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(enc, conf.w, conf.z.Level),
				core,
			)
		}
	}
	if conf.service != "" {
		return zap.New(core, zap.Fields(zap.String("srv", conf.service))), nil
	}
	return zap.New(core), nil
}

// handleOutputFile handles options in log configs to redirect to file
func handleOutputFile(config *zap.Config, file string, fileOnly bool) (zapcore.WriteSyncer, error) {

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
		Compress:   true,
	})

	if fileOnly {
		config.OutputPaths = []string{file}
		return w, nil
	}

	config.OutputPaths = append(config.OutputPaths, file)
	return w, nil
}
