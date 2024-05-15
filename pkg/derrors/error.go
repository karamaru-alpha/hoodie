package derrors

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
	ErrorPattern() ErrorPattern
	SetValues(valueMap map[string]any) Error
}

// derror days's custom error
type derror struct {
	msg      string
	stack    *stack
	pattern  ErrorPattern
	valueMap valueMap
	err      error
}

type valueMap map[string]any

func New(pattern ErrorPattern, msg string) Error {
	return &derror{
		msg:     msg,
		stack:   callers(),
		pattern: pattern,
	}
}

func Wrap(err error, pattern ErrorPattern, msg string) Error {
	var s *stack
	var e *derror
	if errors.As(err, &e) {
		s = caller()
	} else {
		s = callers()
	}

	return &derror{
		msg:     msg,
		stack:   s,
		pattern: pattern,
		err:     err,
	}
}

func As(err error) (Error, bool) {
	var herr *derror
	if errors.As(err, &herr) {
		return herr, true
	}
	return herr, false
}

func (e *derror) Error() string {
	return fmt.Sprintf("%s: %s%s", e.pattern.ErrorCode, e.msg, e.valueMap.String())
}

func (e *derror) Unwrap() error {
	return e.err
}

func (e *derror) Format(s fmt.State, v rune) {
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

func (e *derror) ErrorPattern() ErrorPattern {
	if e != nil {
		return e.pattern
	}
	return Unknown
}

func (e *derror) SetValues(valueMap map[string]any) Error {
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
