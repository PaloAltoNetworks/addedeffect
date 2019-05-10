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
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config for this package
type Config struct {
	z        zap.Config
	w        zapcore.WriteSyncer
	service  string
	fileOnly bool
}

// Option for functional args
type Option func(*Config)

// OptionPrettyTimeStamp prettifies the timestamp
func OptionPrettyTimeStamp() Option {
	return func(config *Config) {
		config.z.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
	}
}

// OptionServiceName configures a service name
func OptionServiceName(name string) Option {
	return func(config *Config) {
		if name != "" {
			config.service = name
		}
	}
}

// OptionFile configures file as logging destination
func OptionFile(file string, fileOnly bool) Option {
	return func(config *Config) {
		if file == "" {
			return
		}

		// Handle log file output
		var err error
		config.fileOnly = fileOnly
		config.w, err = handleOutputFile(&config.z, file, fileOnly)
		if err != nil {
			panic(err)
		}
	}
}

// Configure configures the shared default logger.
func Configure(level string, format string, opts ...Option) zap.Config {

	config := &Config{
		z: getConfig(level, format),
	}

	for _, opt := range opts {
		opt(config)
	}

	logger, err := initLogger(config)
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config.z)

	return config.z
}

// ConfigureWithName configures the shared default logger.
func ConfigureWithName(serviceName string, level string, format string) zap.Config {

	opts := []Option{}
	if serviceName != "" {
		opts = append(opts, OptionServiceName(serviceName))
	}

	return Configure(level, format, opts...)
}

// ConfigureWithOptions configures the shared default logger with options such as file and timestamp formats.
func ConfigureWithOptions(level string, format string, file string, fileOnly bool, prettyTimestamp bool) zap.Config {

	opts := []Option{}
	if file != "" {
		opts = append(opts, OptionFile(file, fileOnly))
	}
	if prettyTimestamp {
		opts = append(opts, OptionPrettyTimeStamp())
	}

	return Configure(level, format, opts...)
}
