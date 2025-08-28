package apperr

import (
	"fmt"
	"runtime"
	"strings"
)

type StackFrame struct {
	Func string
	Line int
	File string
}

func (s *StackFrame) String() string {
	return fmt.Sprintf("%s()\n\t%s:%d", s.Func, s.File, s.Line)
}

type StackTrace struct {
	Frames []StackFrame
}

func NewStackTrace(skip int) *StackTrace {
	var stackFrames []StackFrame
	var pc [100]uintptr
	n := runtime.Callers(skip, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}

		stackFrames = append(stackFrames, StackFrame{
			Func: frame.Function,
			Line: frame.Line,
			File: frame.File,
		})
	}

	return &StackTrace{
		Frames: stackFrames,
	}
}

func (s *StackTrace) String() string {
	var b strings.Builder
	for i, f := range s.Frames {
		b.WriteString(f.String())
		if i < len(s.Frames)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
