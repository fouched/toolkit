package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fouched/toolkit/v2/errors"
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
				frames := errors.Stack(err)
				if len(frames) > 0 {
					fmt.Printf("  %sstack:%s\n", colorCyan, colorReset)
					for _, pc := range frames {
						f := errors.Frame(pc)
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
