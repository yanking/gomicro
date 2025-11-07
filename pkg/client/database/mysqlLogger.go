package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*mysqlLogger)(nil)

type mysqlLogger struct {
	logger *slog.Logger
}

// LogMode 设置日志级别
func (m *mysqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return m
}

// Info 记录 info 级别日志
func (m *mysqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn 记录 warning 级别日志
func (m *mysqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error 记录 error 级别日志
func (m *mysqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace 记录 SQL 执行轨迹
func (m *mysqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		m.logger.Error("MySQL query error",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", elapsed),
			slog.String("error", err.Error()))
		return
	}

	m.logger.Info("MySQL query executed",
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Duration("duration", elapsed))
}
