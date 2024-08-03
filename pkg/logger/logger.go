package logger

import (
	"os"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger/rotator"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ProdLevel = "prod"
)

// A Logger represents an active logging object that uses zap.Logger
// to produce logs.
type Logger struct {
	l *zap.Logger
}

// New creates a new Logger.
func New(rotation bool) *Logger {
	var l Logger
	lvl := os.Getenv("LOG_LEVEL")
	cfg := newZapConfig(lvl)

	core := setupLoggingCore(cfg, rotation)
	l.l = zap.New(core)

	return &l
}

// setupLoggingCore creates a zapcore.Core that duplicates log entries into
// two or more underlying Cores.
func setupLoggingCore(cfg zap.Config, rotation bool) zapcore.Core {
	var cores []zapcore.Core

	if rotation {
		r := rotator.New()
		fSync := zapcore.AddSync(r.Logger)
		fEncoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
		cores = append(cores, zapcore.NewCore(fEncoder, fSync, cfg.Level))
	}

	cSyncer := zapcore.AddSync(zapcore.Lock(os.Stdout))
	cEncoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)

	cores = append(cores, zapcore.NewCore(cEncoder, cSyncer, cfg.Level))

	return zapcore.NewTee(cores...)
}

// getConfig returns a zap.Config of the specified level (ProdLevel or DevLevel).
// The default one is zap.NewDevelopmentConfig.
func newZapConfig(lvl string) zap.Config {
	cfg := zap.NewDevelopmentConfig()

	if lvl == ProdLevel {
		cfg = zap.NewProductionConfig()
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return cfg
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.l.Debug(msg, fields...)
}

// Info logs a message at Info level.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

// Warn logs a message at Warn level.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.l.Warn(msg, fields...)
}

// Error logs a message at Error level.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.l.Error(msg, fields...)
}

// Fatal logs a message at Fatal level.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.l.Fatal(msg, fields...)
}
