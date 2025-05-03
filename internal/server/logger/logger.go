package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var ZapLogger = zap.NewNop()

func Initialize(level string) (*zap.Logger, error) {
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

	ZapLogger, err = cfg.Build()

	return ZapLogger, err
}
