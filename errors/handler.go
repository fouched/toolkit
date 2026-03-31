package errors

import (
	"context"
	"log/slog"
)

type ErrorHandler struct {
	slog.Handler
}

func (h *ErrorHandler) Handle(ctx context.Context, r slog.Record) error {
	// Walk through all attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			if err, ok := a.Value.Any().(error); ok {
				// Replace the "err" field with the error string
				r.Add("err", slog.StringValue(err.Error()))

				// Add the stack trace automatically
				r.Add("stack", slog.AnyValue(Stack(err)))
			}
		}
		return true
	})

	return h.Handler.Handle(ctx, r)
}
