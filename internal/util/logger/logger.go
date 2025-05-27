package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(levelStr string) (*zap.Logger, error) {
	lvl := zapcore.InfoLevel
	if err := lvl.UnmarshalText([]byte(levelStr)); err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(lvl)
	return cfg.Build()
}
