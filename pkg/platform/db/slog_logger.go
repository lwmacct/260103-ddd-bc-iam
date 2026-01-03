package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SlogLogger 实现 GORM logger.Interface，将日志输出到 slog
type SlogLogger struct {
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewSlogLogger 创建一个新的 slog logger
func NewSlogLogger(level logger.LogLevel) *SlogLogger {
	return &SlogLogger{
		LogLevel:                  level,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 设置日志级别
func (l *SlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 输出 info 级别日志
func (l *SlogLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		pc := logm.CallerPC("gorm.io/gorm", "pkg/platform/db")
		logm.LogWithPC(ctx, slog.LevelInfo, pc, fmt.Sprintf(msg, data...))
	}
}

// Warn 输出 warn 级别日志
func (l *SlogLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		pc := logm.CallerPC("gorm.io/gorm", "pkg/platform/db")
		logm.LogWithPC(ctx, slog.LevelWarn, pc, fmt.Sprintf(msg, data...))
	}
}

// Error 输出 error 级别日志
func (l *SlogLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		pc := logm.CallerPC("gorm.io/gorm", "pkg/platform/db")
		logm.LogWithPC(ctx, slog.LevelError, pc, fmt.Sprintf(msg, data...))
	}
}

// Trace 输出 SQL 追踪日志
func (l *SlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取真正的调用者 PC，跳过 GORM 和 database 包内部调用
	pc := logm.CallerPC("gorm.io/gorm", "infrastructure/database")

	// 构建日志属性
	attrs := []any{
		slog.Duration("elapsed", elapsed),
		slog.Int64("rows", rows),
		slog.String("sql", sql),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		// 忽略 record not found 错误
		if l.IgnoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		logm.LogWithPC(ctx, slog.LevelError, pc, "GORM query error", append(attrs, slog.Any("error", err))...)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		logm.LogWithPC(ctx, slog.LevelWarn, pc, "GORM slow query", append(attrs, slog.Duration("threshold", l.SlowThreshold))...)

	case l.LogLevel >= logger.Info:
		logm.LogWithPC(ctx, slog.LevelDebug, pc, "GORM query", attrs...)
	}
}
