package logging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fouched/toolkit/v2/faults"
)

// ProdHandler is the production logging handler.
// It ensures errors are logged with message + stack frames,
// while keeping all other attributes structured and clean.
type ProdHandler struct {
	slog.Handler
}

var ProdLevel = new(slog.LevelVar)

func (h *ProdHandler) Handle(ctx context.Context, r slog.Record) error {
	// Build a new record with the same metadata
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	// Copy all attributes, but intercept "err"
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			if err, ok := a.Value.Any().(error); ok {
				// Use ONLY the message chain
				newRecord.Add("err", err.Error())

				// Add formatted stack frames
				frames := faults.Stack(err)
				formatted := make([]string, len(frames))
				for i, pc := range frames {
					f := faults.Frame(pc)
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
