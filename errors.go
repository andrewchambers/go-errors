package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

type Frame struct {
	Message string
	PCKnown bool
	Values  []KV
	PC      uintptr
}

type Trace struct {
	Frames []Frame
}

type KV struct {
	K string
	V interface{}
}

type ErrorContext struct {
	Values []KV
}

type Error struct {
	// A message for th eend user.
	Message string
	// Error context, intended to be immutable,
	// Never mutate values stored in the context.
	// Never modify this in place.
	Values []KV
	// Program counter of where this error originates.
	SourcePC uintptr
	// The original cause of the error, if nil,
	// the error itself is the cause,
	Cause error
	// Root cause of the error chain.
	RootCause error
}

func (err *Error) String() string {
	return err.Error()
}

func (err *Error) Error() string {
	rootCause := RootCause(err)

	if rootCause == nil || err.Cause == nil {
		return err.Message
	}

	return fmt.Sprintf("%s: %s", err.Message, rootCause.Error())
}

func New(msg string) error {
	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   msg,
		SourcePC:  pc[0],
		Values:    nil,
		Cause:     nil,
		RootCause: nil,
	}
}

func Errorf(format string, args ...interface{}) error {

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   fmt.Sprintf(format, args...),
		SourcePC:  pc[0],
		Values:    nil,
		Cause:     nil,
		RootCause: nil,
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   msg,
		SourcePC:  pc[0],
		Values:    nil,
		Cause:     err,
		RootCause: RootCause(err),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   fmt.Sprintf(format, args...),
		SourcePC:  pc[0],
		Values:    nil,
		Cause:     err,
		RootCause: RootCause(err),
	}
}

func Context(values []KV) ErrorContext {
	return ErrorContext{Values: values}
}

func (c ErrorContext) Context(values []KV) ErrorContext {
	return ErrorContext{Values: values}
}

func (c ErrorContext) New(msg string) error {
	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   msg,
		SourcePC:  pc[0],
		Values:    c.Values,
		Cause:     nil,
		RootCause: nil,
	}
}

func (c ErrorContext) Errorf(format string, args ...interface{}) error {

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   fmt.Sprintf(format, args...),
		SourcePC:  pc[0],
		Values:    c.Values,
		Cause:     nil,
		RootCause: nil,
	}
}

func (c ErrorContext) Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   msg,
		SourcePC:  pc[0],
		Values:    c.Values,
		Cause:     err,
		RootCause: RootCause(err),
	}
}

func (c ErrorContext) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var pc [1]uintptr
	runtime.Callers(2, pc[:])

	return &Error{
		Message:   fmt.Sprintf(format, args...),
		SourcePC:  pc[0],
		Values:    c.Values,
		Cause:     err,
		RootCause: RootCause(err),
	}
}

// Return an error trace up to the maximum of 10000 frames.
func GetTrace(err error) *Trace {
	t := &Trace{}

	for i := 0; i < 10000; i++ {
		if err == nil {
			return t
		}
		e, ok := err.(*Error)
		if !ok {
			t.Frames = append(t.Frames, Frame{
				Message: err.Error(),
				PCKnown: false,
			})
			return t
		}

		t.Frames = append(t.Frames, Frame{
			Message: e.Message,
			PCKnown: true,
			PC:      e.SourcePC,
			Values:  e.Values,
		})
		err = e.Cause
	}

	return t
}

// Return the original error cause of an error if possible.
func RootCause(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if ok {
		if e.RootCause != nil {
			return e.RootCause
		}
	}

	return err
}

func (t *Trace) String() string {
	var buf bytes.Buffer

	hasPrev := false
	for _, f := range t.Frames {
		if hasPrev {
			_, _ = fmt.Fprintf(&buf, "Cause:\n")
		}
		hasPrev = true
		if f.PCKnown {
			fn := runtime.FuncForPC(f.PC)
			file, line := fn.FileLine(f.PC)
			_, _ = fmt.Fprintf(&buf, "%s:%s:%d %#v\n", file, fn.Name(), line, f.Message)
			_, _ = fmt.Fprintf(&buf, "Where:\n")
			for _, kv := range f.Values {
				_, _ = fmt.Fprintf(&buf, "  %#v = %#v\n", kv.K, kv.V)
			}
		} else {
			_, _ = fmt.Fprintf(&buf, "?:? %#v\n", f.Message)
		}
	}

	return buf.String()
}
