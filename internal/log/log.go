package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	L *zap.Logger
	S *zap.SugaredLogger
)

func Init(verbose bool) {
	level := zapcore.InfoLevel
	if verbose {
		level = zapcore.DebugLevel
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stderr),
		level,
	)

	L = zap.New(core)
	S = L.Sugar()
}

func Sync() {
	if L != nil {
		_ = L.Sync()
	}
}
