package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

// Logger wraps log.Logger and implements util.Logger interface
type Logger struct {
	*log.Logger
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
	logger := slog.NewLogLogger(logHandler, level.Level())
	slog.SetDefault(slogger)

	return &Logger{Logger: logger}, nil
}

// Infof implements util.Logger interface
func (l *Logger) Infof(format string, v ...any) {
	slog.Info(fmt.Sprintf(format, v...))
}

// Errorf implements util.Logger interface
func (l *Logger) Errorf(format string, v ...any) {
	slog.Error(fmt.Sprintf(format, v...))
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
