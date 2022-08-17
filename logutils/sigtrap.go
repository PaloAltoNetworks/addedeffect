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

//go:build linux || darwin
// +build linux darwin

package logutils

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func handleElevationSignal(cfg zap.Config) {

	defaultLevel := cfg.Level
	var elevated bool

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for s := range c {
		if s == syscall.SIGINT {
			return
		}
		elevated = !elevated

		if elevated {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			l, _ := cfg.Build()
			zap.ReplaceGlobals(l)
			zap.L().Info("Log level elevated to debug")
		} else {
			zap.L().Info("Log level restored to original configuration", zap.Stringer("level", defaultLevel.Level()))
			cfg.Level = defaultLevel
			l, _ := cfg.Build()
			zap.ReplaceGlobals(l)
		}
	}
}
