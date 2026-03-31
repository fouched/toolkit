package errors

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	msg   string
	cause error
	stack []uintptr
}

func New(msg string) error {
	return &Error{
		msg:   msg,
		stack: callers(),
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	// If the cause already has a stack, preserve it
	var e *Error
	if errors.As(err, &e) {
		return &Error{
			msg:   msg,
			cause: err,
			stack: e.stack, // <-- preserve deepest stack
		}
	}

	// Otherwise capture a new stack
	return &Error{
		msg:   msg,
		cause: err,
		stack: callers(),
	}
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		msg:   err.Error(),
		cause: err,
		stack: callers(),
	}
}

func (e *Error) Error() string {
	if e.cause == nil {
		return e.msg
	}
	return e.msg + ": " + e.cause.Error()
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// Print the full message chain
			fmt.Fprint(s, e.Error())

			// Print only THIS error's stack
			for _, pc := range e.stack {
				f := Frame(pc)
				fmt.Fprintf(s, "\n  at %+v", f)
			}
			return
		}
		fallthrough
	default:
		fmt.Fprint(s, e.Error())
	}
}

func callers() []uintptr {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(3, pcs)
	return pcs[:n]
}
