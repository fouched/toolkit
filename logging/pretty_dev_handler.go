package logging

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/fouched/toolkit/v2/faults"
)

type PrettyDevHandler struct{}

func NewPrettyDevHandler() *PrettyDevHandler {
	return &PrettyDevHandler{}
}

func (h *PrettyDevHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *PrettyDevHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := colorForLevel(r.Level)

	// Determine caller frame for prefix
	prefix := ""
	if frame, ok := callerFrame(); ok {
		prefix = prefixForFrame(frame)
	}

	// HEADER
	fmt.Printf("%s%s%s %s %s%s\n",
		levelColor,
		r.Level.String(),
		colorReset,
		r.Time.Format(time.RFC3339),
		prefix,
		r.Message,
	)

	// ATTRIBUTES
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "err":
			if err, ok := a.Value.Any().(error); ok {
				fmt.Printf("  %serr:%s %s\n", colorRed, colorReset, err.Error())

				// Stack
				frames := faults.Stack(err)
				if len(frames) > 0 {
					fmt.Printf("  %sstack:%s\n", colorCyan, colorReset)
					for _, pc := range frames {
						f := faults.Frame(pc)
						fmt.Printf("    %s:%d  %s\n", f.File(), f.Line(), f.Function())
					}
				}
				return true
			}

		default:
			fmt.Printf("  %s=%v\n", a.Key, a.Value.Any())
		}
		return true
	})

	fmt.Println()
	return nil
}

func (h *PrettyDevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyDevHandler) WithGroup(name string) slog.Handler {
	return h
}

// ------------------------------------------------------------
// Prefix logic
// ------------------------------------------------------------

func callerFrame() (runtime.Frame, bool) {
	pcs := make([]uintptr, 1)
	// Skip slog internals + handler
	n := runtime.Callers(5, pcs)
	if n == 0 {
		return runtime.Frame{}, false
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	return frame, true
}

func prefixForFrame(f runtime.Frame) string {
	path := f.File

	// List of known layer directories
	layers := []string{
		"handlers",
		"services",
		"service",
		"repo",
		"repository",
		"repositories",
		"store",
		"persistence",
		"domain",
		"usecase",
		"http",
		"api",
	}

	for _, layer := range layers {
		needle := "/" + layer + "/"
		if strings.Contains(path, needle) {
			return colorForLayer(layer) + layer + " → " + colorReset
		}
	}

	return ""
}

func colorForLayer(layer string) string {
	switch layer {
	case "handlers", "http", "api":
		return colorBlue
	case "services", "service", "domain", "usecase":
		return colorCyan
	case "repo", "repository", "repositories", "store", "persistence":
		return colorGreen
	default:
		return colorReset
	}
}
