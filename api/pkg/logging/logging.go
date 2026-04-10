package logging

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const loggerKey = contextKey("logger")

type Logger struct {
	logger *slog.Logger
	prefix string
}

func New() *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

// WithField returns a new Logger that includes a structured key-value field
// on every log entry. Can be chained: logger.WithField("service", "cookbooks").WithField("version", "2")
func (l *Logger) WithField(key, value string) *Logger {
	return &Logger{
		logger: l.logger.With(key, value),
		prefix: l.prefix,
	}
}

// WithPrefix returns a new Logger that prepends a string prefix to every message.
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		logger: l.logger,
		prefix: prefix,
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Info(l.prefix+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logger.Debug(l.prefix+format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.Warn(l.prefix+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Error(l.prefix+format, v...)
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return New()
}
