package logging

import "log/slog"

type layerDef struct {
	name   string
	needle string
}

var layers = []layerDef{
	{"handlers", "/handlers/"},
	{"services", "/services/"},
	{"service", "/service/"},
	{"repo", "/repo/"},
	{"repository", "/repository/"},
	{"repositories", "/repositories/"},
	{"store", "/store/"},
	{"persistence", "/persistence/"},
	{"domain", "/domain/"},
	{"usecase", "/usecase/"},
	{"http", "/http/"},
	{"api", "/api/"},
}

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
)

func colorForLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return colorMagenta
	case slog.LevelInfo:
		return colorGreen
	case slog.LevelWarn:
		return colorYellow
	case slog.LevelError:
		return colorRed
	default:
		return colorReset
	}
}
