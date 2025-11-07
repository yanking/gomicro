// Package database provides support for multiple database instances including MySQL, Redis and MongoDB.
package database

import (
	"fmt"
	"log/slog"
)

// mongoLogger is a logger for MongoDB operations
type mongoLogger struct {
	logger *slog.Logger
}

// newMongoLogger creates a new MongoDB logger
func newMongoLogger(logger *slog.Logger) *mongoLogger {
	return &mongoLogger{
		logger: logger,
	}
}

// Info logs info level messages
func (l *mongoLogger) Info(_ int, msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValuesToAttr(keysAndValues)...)
}

// Error logs error level messages
func (l *mongoLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	attrs := keysAndValuesToAttr(keysAndValues)
	attrs = append(attrs, slog.String("error", err.Error()))
	l.logger.Error(msg, attrs...)
}

// keysAndValuesToAttr converts keysAndValues to slog.Attr
func keysAndValuesToAttr(keysAndValues []interface{}) []any {
	if len(keysAndValues) == 0 {
		return nil
	}

	attrs := make([]any, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			attrs = append(attrs, slog.Any(fmt.Sprintf("%v", keysAndValues[i]), keysAndValues[i+1]))
		}
	}
	return attrs
}
