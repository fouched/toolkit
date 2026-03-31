package errors

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

type ErrorHandler struct {
	slog.Handler
}

type PrettyDevHandler struct{}

func (h *ErrorHandler) Handle(ctx context.Context, r slog.Record) error {
	// Build a new record with the same metadata
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	// Copy all attributes, but intercept "err"
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			if err, ok := a.Value.Any().(error); ok {
				// Use ONLY the message chain
				newRecord.Add("err", err.Error())

				// Add formatted stack frames
				frames := Stack(err)
				formatted := make([]string, len(frames))
				for i, pc := range frames {
					f := Frame(pc)
					formatted[i] = fmt.Sprintf("%+v", f)
				}
				newRecord.Add("stack", formatted)
				return true
			}
		}

		// Keep all other attributes unchanged
		newRecord.Add(a.Key, a.Value)
		return true
	})

	return h.Handler.Handle(ctx, newRecord)
}

func NewPrettyDevHandler() *PrettyDevHandler {
	return &PrettyDevHandler{}
}

func (h *PrettyDevHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}
func (h *PrettyDevHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := colorForLevel(r.Level)

	// HEADER
	fmt.Printf("%s%s%s %s %s\n",
		levelColor,
		r.Level.String(),
		colorReset,
		r.Time.Format(time.RFC3339),
		r.Message,
	)

	// ATTRIBUTES
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "err":
			if err, ok := a.Value.Any().(error); ok {
				fmt.Printf("  %serr:%s %s\n", colorRed, colorReset, err.Error())

				// Stack
				frames := Stack(err)
				if len(frames) > 0 {
					fmt.Printf("  %sstack:%s\n", colorCyan, colorReset)
					for _, pc := range frames {
						f := Frame(pc)
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

func (h *PrettyDevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyDevHandler) WithGroup(name string) slog.Handler {
	return h
}
