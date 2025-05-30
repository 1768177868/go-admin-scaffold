package logger

import (
	"go.uber.org/zap"
)

// Logger 封装了zap.Logger
type Logger struct {
	logger *zap.Logger
}

// Info 记录info级别日志
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if len(fields) > 0 {
		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
		l.logger.Info(msg, zapFields...)
	} else {
		l.logger.Info(msg)
	}
}

// Error 记录error级别日志
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if len(fields) > 0 {
		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
		l.logger.Error(msg, zapFields...)
	} else {
		l.logger.Error(msg)
	}
}

// Debug 记录debug级别日志
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if len(fields) > 0 {
		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
		l.logger.Debug(msg, zapFields...)
	} else {
		l.logger.Debug(msg)
	}
}

// Warn 记录warn级别日志
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	if len(fields) > 0 {
		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
		l.logger.Warn(msg, zapFields...)
	} else {
		l.logger.Warn(msg)
	}
}

// Close 关闭日志
func (l *Logger) Close() error {
	return l.logger.Sync()
}
