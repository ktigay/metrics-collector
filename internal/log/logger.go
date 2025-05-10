package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AppLogger = zap.NewNop()
var SugaredLogger = *AppLogger.Sugar()

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

	AppLogger, err = cfg.Build()
	SugaredLogger = *AppLogger.Sugar()

	return AppLogger, err
}
