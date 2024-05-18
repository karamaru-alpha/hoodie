package perrors

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
)

type Error interface {
	error
	fmt.Formatter
	Unwrap() error
	SetValues(valueMap map[string]any) Error
}

// perror days's custom error
type perror struct {
	msg      string
	stack    *stack
	valueMap valueMap
	err      error
}

type valueMap map[string]any

func New(msg string) Error {
	return &perror{
		msg:   msg,
		stack: callers(),
	}
}

func Wrap(err error, msg string) Error {
	var s *stack
	var e *perror
	if errors.As(err, &e) {
		s = caller()
	} else {
		s = callers()
	}

	return &perror{
		msg:   msg,
		stack: s,
		err:   err,
	}
}

func As(err error) (Error, bool) {
	var herr *perror
	if errors.As(err, &herr) {
		return herr, true
	}
	return herr, false
}

func (e *perror) Error() string {
	return fmt.Sprintf("%s%s", e.msg, e.valueMap.String())
}

func (e *perror) Unwrap() error {
	return e.err
}

func (e *perror) Format(s fmt.State, v rune) {
	switch v {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%s\n%s", e.Error(), strings.Join(e.stack.GetStackTrace(), "\n"))
			if e.Unwrap() != nil {
				_, _ = fmt.Fprintf(s, "\n- %+v", e.Unwrap())
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *perror) SetValues(valueMap map[string]any) Error {
	e.valueMap = valueMap
	return e
}

func (v valueMap) String() string {
	if len(v) == 0 {
		return ""
	}

	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ss := make([]string, 0, len(keys)+1)
	ss = append(ss, "")
	for _, key := range keys {
		ss = append(ss, fmt.Sprintf("%s=%+v", key, v[key]))
	}
	return strings.Join(ss, " ")
}

type stack []uintptr

func (s *stack) GetStackTrace() []string {
	if s == nil {
		return nil
	}
	frames := runtime.CallersFrames(*s)
	stackTrace := make([]string, 0)
	for {
		frame, more := frames.Next()
		stackTrace = append(stackTrace, fmt.Sprintf("\t%s\n\t\t%s:%d", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return stackTrace
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func caller() *stack {
	pc, _, _, _ := runtime.Caller(2)
	return &stack{pc}
}
