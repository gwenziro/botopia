package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
)

// LogLevel mendefinisikan level logging
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger adalah wrapper sederhana untuk logging
type Logger struct {
	Module    string
	Level     LogLevel
	UseColors bool
	writer    io.Writer
}

// New membuat instance logger baru
func New(module string, level LogLevel, useColors bool) *Logger {
	writer := colorable.NewColorableStdout()
	return &Logger{
		Module:    module,
		Level:     level,
		UseColors: useColors,
		writer:    writer,
	}
}

// SetWriter mengatur writer untuk output log
func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

// LevelFromString mengkonversi string level menjadi LogLevel
func LevelFromString(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// levelColor mengembalikan kode warna ANSI untuk level
func (l *Logger) levelColor(level LogLevel) string {
	if !l.UseColors {
		return ""
	}

	switch level {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// levelName mengembalikan nama level
func levelName(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// resetColor mengembalikan kode reset warna ANSI
func (l *Logger) resetColor() string {
	if !l.UseColors {
		return ""
	}
	return "\033[0m"
}

// log adalah fungsi internal untuk logging
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.Level {
		return
	}

	now := time.Now().Format("15:04:05.000")
	levelStr := levelName(level)
	color := l.levelColor(level)
	reset := l.resetColor()

	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("%s%s [%s] %s: %s%s\n",
		color, now, levelStr, l.Module, message, reset)

	fmt.Fprint(l.writer, logLine)

	// Flush jika fatal
	if level == FATAL {
		if f, ok := l.writer.(interface{ Flush() error }); ok {
			f.Flush()
		}
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}
