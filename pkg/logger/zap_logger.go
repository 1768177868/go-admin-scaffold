package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

type Config struct {
	Level      string `yaml:"level"`       // 日志级别
	Filename   string `yaml:"filename"`    // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 每个日志文件最大尺寸，单位MB
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"max_age"`     // 保留的旧日志文件最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧日志文件
	Daily      bool   `yaml:"daily"`       // 是否按天切割日志
}

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Setup initializes the logger
func Setup(config *Config) error {
	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create log writer
	var writer zapcore.WriteSyncer
	if config.Daily {
		writer = getDailyWriter(config)
	} else {
		writer = getWriter(config)
	}

	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		level,
	)

	// Create logger
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()

	return nil
}

// getWriter creates a lumberjack writer for continuous log file
func getWriter(config *Config) zapcore.WriteSyncer {
	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(config.Filename), 0744); err != nil {
		panic(err)
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})
}

// getDailyWriter creates a writer that rotates log files daily
func getDailyWriter(config *Config) zapcore.WriteSyncer {
	// Ensure log directory exists
	logDir := filepath.Dir(config.Filename)
	if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(err)
	}

	// Get base filename without extension
	base := filepath.Base(config.Filename)
	ext := filepath.Ext(base)
	prefix := base[:len(base)-len(ext)]

	// Create daily log file name
	dailyFile := filepath.Join(logDir, fmt.Sprintf("%s-%s%s",
		prefix,
		time.Now().Format("2006-01-02"),
		ext,
	))

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   dailyFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})
}

// WithField adds a field to the logger context
func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

// getTraceID gets trace ID from context
func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, args ...interface{}) {
	sugar.Debugw(msg, append(args, "trace_id", getTraceID(ctx))...)
}

// Info logs an info message
func Info(ctx context.Context, msg string, args ...interface{}) {
	sugar.Infow(msg, append(args, "trace_id", getTraceID(ctx))...)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, args ...interface{}) {
	sugar.Warnw(msg, append(args, "trace_id", getTraceID(ctx))...)
}

// Error logs an error message
func Error(ctx context.Context, msg string, args ...interface{}) {
	sugar.Errorw(msg, append(args, "trace_id", getTraceID(ctx))...)
}

// Fatal logs a fatal message and exits
func Fatal(ctx context.Context, msg string, args ...interface{}) {
	sugar.Fatalw(msg, append(args, "trace_id", getTraceID(ctx))...)
}

// Close flushes any buffered log entries
func Close() error {
	return logger.Sync()
}
