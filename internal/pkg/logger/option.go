package logger

import "go.uber.org/zap"

type ConfigOption func(c *zap.Config)

func WithDebugLevel() ConfigOption {
	return func(c *zap.Config) {
		c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}

func WithOutputStdout() ConfigOption {
	return func(c *zap.Config) {
		c.OutputPaths = []string{"stdout"}
	}
}
