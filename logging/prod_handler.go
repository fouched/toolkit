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
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	// Copy all attributes
	r.Attrs(func(a slog.Attr) bool {
		val := a.Value.Any()

		// Intercept ANY error, regardless of key
		if err, ok := val.(error); ok {
			wrapped := faults.Wrap(err, r.Message)

			// Add message chain
			newRecord.Add("err", wrapped.Error())

			// Add stack frames
			frames := faults.Stack(wrapped)
			formatted := make([]string, len(frames))
			for i, pc := range frames {
				f := faults.Frame(pc)
				formatted[i] = fmt.Sprintf("%+v", f)
			}
			newRecord.Add("stack", formatted)

			return true
		}

		// Normal attribute
		newRecord.Add(a.Key, a.Value)
		return true
	})

	return h.Handler.Handle(ctx, newRecord)
}
