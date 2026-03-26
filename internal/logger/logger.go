package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 構造化ロガー
type Logger struct {
	*zap.SugaredLogger
}

// New 新しいロガーを作成
func New(level string) *Logger {
	// レベルのパース
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// コンフィグ作成
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// ロガー作成
	logger, err := config.Build()
	if err != nil {
		panic("failed to build logger: " + err.Error())
	}
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// Debug デバッグログ
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.SugaredLogger.Debugw(msg, keyvals...)
}

// Info 情報ログ
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.SugaredLogger.Infow(msg, keyvals...)
}

// Warn 警告ログ
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.SugaredLogger.Warnw(msg, keyvals...)
}

// Error エラーログ
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.SugaredLogger.Errorw(msg, keyvals...)
}

// Sync ログをフラッシュ
func (l *Logger) Sync() {
	l.SugaredLogger.Sync()
}
