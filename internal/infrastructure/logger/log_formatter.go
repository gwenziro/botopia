package logger

import (
	"fmt"
	"time"
)

// LogFormat mengatur format pesan log
type LogFormat struct {
	TimeFormat  string
	ShowModule  bool
	ColorOutput bool
}

// DefaultLogFormat adalah format log default
var DefaultLogFormat = LogFormat{
	TimeFormat:  "15:04:05.000",
	ShowModule:  true,
	ColorOutput: true,
}

// FormatLog memformat pesan log
func FormatLog(level LogLevel, module, message string, format LogFormat) string {
	now := time.Now().Format(format.TimeFormat)
	levelStr := levelName(level)

	var color, reset string
	if format.ColorOutput {
		color = levelColor(level)
		reset = "\033[0m"
	}

	var moduleInfo string
	if format.ShowModule && module != "" {
		moduleInfo = fmt.Sprintf(" %s:", module)
	}

	return fmt.Sprintf("%s%s [%s]%s %s%s",
		color, now, levelStr, moduleInfo, message, reset)
}

// levelColor mengembalikan kode warna ANSI untuk level
func levelColor(level LogLevel) string {
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
