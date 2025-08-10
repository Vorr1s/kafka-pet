package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func New(level zapcore.LevelEnabler, w io.Writer, opt ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	enc := zapcore.NewConsoleEncoder(cfg)
	return zap.New(zapcore.NewCore(enc, zapcore.AddSync(w), level)).Sugar()
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func GetLogger() *zap.SugaredLogger {
	return global
}

func Init() {
	SetLogger(New(defaultLevel, os.Stdout))
}
