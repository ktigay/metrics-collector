// Package log Логгер.
package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initialize инициализация глобального логгера.
func Initialize(level string) (*zap.SugaredLogger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig = encoderConfig
	cfg.Level = lvl

	var logger *zap.Logger
	if logger, err = cfg.Build(); err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
