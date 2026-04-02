package faults

import (
	"errors"
	"fmt"
	"runtime"
)

type Frame uintptr

func (f Frame) pc() uintptr {
	return uintptr(f)
}

func (f Frame) File() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) Function() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func Stack(err error) []Frame {
	var e *Error
	if errors.As(err, &e) {
		frames := make([]Frame, len(e.stack))
		for i, pc := range e.stack {
			frames[i] = Frame(pc)
		}
		return frames
	}
	return nil
}

// Format implements fmt.Formatter.
// Supported verbs:
//
//	%v  - short file:line
//	%+v - full file:line with function name
//	%s  - file:line
//	%n  - function name only
func (f Frame) Format(s fmt.State, verb rune) {
	pc := uintptr(f)
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Fprint(s, "unknown")
		return
	}

	file, line := fn.FileLine(pc)

	switch verb {
	case 'v':
		if s.Flag('+') {
			// Single-line full format: func (file:line)
			fmt.Fprintf(s, "%s (%s:%d)", fn.Name(), file, line)
		} else {
			// Short file:line
			short := shortFile(file)
			fmt.Fprintf(s, "%s:%d", short, line)
		}

	case 's':
		fmt.Fprintf(s, "%s:%d", file, line)

	case 'n':
		fmt.Fprint(s, fn.Name())

	default:
		fmt.Fprintf(s, "%s:%d", file, line)
	}
}

func shortFile(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
