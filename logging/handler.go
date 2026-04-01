package logging

import (
	"log/slog"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func colorForLevel(level slog.Level) string {
	switch level {
	case slog.LevelError:
		return colorRed
	case slog.LevelWarn:
		return colorYellow
	case slog.LevelInfo:
		return colorBlue
	default:
		return colorCyan
	}
}
