package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

// Logger wraps slog.Logger and provides a log.Logger adapter for APIs that require it
type Logger struct {
	slogLogger *slog.Logger
	stdLogger  *log.Logger
}

// New creates a new Logger based on the provided configuration
func New(cfg *config.Config) (*Logger, error) {
	level := &slog.LevelVar{}
	level.Set(parseLevel(cfg.LogLevel()))
	logWriter := os.Stderr
	log.SetOutput(logWriter)

	var logHandler slog.Handler
	switch cfg.LogFormat() {
	case "text":
		logHandler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: level})
	case "json":
		logHandler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: level})
	default:
		return nil, fmt.Errorf("unknown log format: %s", cfg.LogFormat())
	}

	slogger := slog.New(logHandler)
	slog.SetDefault(slogger)

	return &Logger{
		slogLogger: slogger,
		stdLogger:  slog.NewLogLogger(logHandler, slog.LevelError),
	}, nil
}

// SlogLogger returns the underlying slog.Logger
func (l *Logger) SlogLogger() *slog.Logger {
	return l.slogLogger
}

// StdLogger returns a log.Logger that writes to the underlying slog handler
func (l *Logger) StdLogger() *log.Logger {
	return l.stdLogger
}

// parseLevel converts string level to slog.Level
func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
