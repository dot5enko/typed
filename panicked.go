package typed

import (
	"fmt"
	"runtime"
)

func DebugStack(lines ...int) []string {
	stack := []string{}

	callers := 10

	if len(lines) > 0 {
		callers = lines[0]
	}

	for i := 0; i < callers; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			stack = append(stack, fmt.Sprintf(" %s : %d", file, line))
		} else {
			break
		}
	}

	return stack
}

var CaptureStackOnPanic = true
var CaptureStackDepth = 20

type PanickedError struct {
	Cause any
	Stack []string
}

func (pe PanickedError) Error() string {
	return fmt.Sprintf("function panicked: %v", pe.Cause)
}

func RecoverIfPanic[T any](resultReference *Result[T]) {
	rec := recover()
	if rec != nil {

		stack := []string{}

		if CaptureStackOnPanic {
			stack = DebugStack(CaptureStackDepth)
		}

		*resultReference = ResultFailed[T](PanickedError{
			Cause: rec,
			Stack: stack,
		})
	}
}
