package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ProdLevel = "prod"
	DevLevel  = "dev"
)

type Logger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

// New creates a new Logger instance.
func New() *Logger {
	var config zap.Config
	level := os.Getenv("LOG_LEVEL")

	switch level {
	case ProdLevel:
		config = zap.NewProductionConfig()
	case DevLevel:
		config = zap.NewDevelopmentConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Lumberjack logger for file rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "/var/log/app/app.log",
		MaxSize:    5,  // megabytes
		MaxBackups: 10, // number of backups
		MaxAge:     14, // days
		Compress:   true,
	}

	// File writer
	fileSyncer := zapcore.AddSync(lumberjackLogger)
	consoleSyncer := zapcore.AddSync(zapcore.Lock(os.Stdout))

	// Configure the encoder
	encoderConfig := config.EncoderConfig
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleSyncer, config.Level),
		zapcore.NewCore(fileEncoder, fileSyncer, config.Level),
	)

	zapLogger := zap.New(core)

	return &Logger{
		logger:        zapLogger,
		sugaredLogger: zapLogger.Sugar(),
	}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() {
	err := l.logger.Sync()
	if err != nil {
		l.logger.Error("failed to sync logger", zap.Error(err))
	}
	err = l.sugaredLogger.Sync()
	if err != nil {
		l.logger.Error("failed to sync logger", zap.Error(err))
	}
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs a message at Info level.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a message at Warn level.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs a message at Error level.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// DPanic logs a message at DPanic level.
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// Panic logs a message at Panic level.
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal logs a message at Fatal level.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugf logs a message at Debug level with a template.
func (l *Logger) Debugf(template string, args ...any) {
	l.sugaredLogger.Debugf(template, args...)
}

// Infof logs a message at Info level with a template.
func (l *Logger) Infof(template string, args ...any) {
	l.sugaredLogger.Infof(template, args...)
}

// Warnf logs a message at Warn level with a template.
func (l *Logger) Warnf(template string, args ...any) {
	l.sugaredLogger.Warnf(template, args...)
}

// Errorf logs a message at Error level with a template.
func (l *Logger) Errorf(template string, args ...any) {
	l.sugaredLogger.Errorf(template, args...)
}

// DPanicf logs a message at DPanic level with a template.
func (l *Logger) DPanicf(template string, args ...any) {
	l.sugaredLogger.DPanicf(template, args...)
}

// Panicf logs a message at Panic level with a template.
func (l *Logger) Panicf(template string, args ...any) {
	l.sugaredLogger.Panicf(template, args...)
}

// Fatalf logs a message at Fatal level with a template.
func (l *Logger) Fatalf(template string, args ...any) {
	l.sugaredLogger.Fatalf(template, args...)
}
