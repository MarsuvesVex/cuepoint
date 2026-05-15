package bot

import (
	"log"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

func NewLogger(level string, logger *log.Logger) *Logger {
	if logger == nil {
		logger = log.Default()
	}
	return &Logger{
		level:  parseLogLevel(level),
		logger: logger,
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logf(LogLevelDebug, "DEBUG", format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logf(LogLevelInfo, "INFO", format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logf(LogLevelWarn, "WARN", format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logf(LogLevelError, "ERROR", format, args...)
}

func (l *Logger) logf(level LogLevel, label, format string, args ...any) {
	if l == nil || l.logger == nil || level < l.level {
		return
	}
	l.logger.Printf("bot[%s] "+format, append([]any{label}, args...)...)
}

func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LogLevelDebug
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "info", "":
		return LogLevelInfo
	default:
		return LogLevelInfo
	}
}
