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
	signal.Notify(c, os.Interrupt, syscall.SIGUSR1)
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
