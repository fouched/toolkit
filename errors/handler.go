package errors

import (
	"context"
	"fmt"
	"log/slog"
)

type ErrorHandler struct {
	slog.Handler
}

type PrettyDevHandler struct {
	slog.Handler
}

func (h *ErrorHandler) Handle(ctx context.Context, r slog.Record) error {
	// Build a new record with the same metadata
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	// Copy all attributes, but intercept "err"
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			if err, ok := a.Value.Any().(error); ok {
				// Replace "err" with the error string
				newRecord.Add("err", err.Error())

				// Add formatted stack frames
				frames := Stack(err)
				formatted := make([]string, len(frames))
				for i, pc := range frames {
					f := pc
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

func (h *PrettyDevHandler) Handle(ctx context.Context, r slog.Record) error {
	// Print the base record first
	err := h.Handler.Handle(ctx, r)
	if err != nil {
		return err
	}

	// Now pretty-print the stack if present
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "stack" {
			if frames, ok := a.Value.Any().([]string); ok {
				fmt.Println("STACK TRACE:")
				for _, f := range frames {
					fmt.Println("  ", f)
				}
			}
		}
		return true
	})

	return nil
}
