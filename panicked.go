package typed

import (
	"fmt"
	"runtime/debug"
)

var CaptureStackOnPanic = true

type PanickedError struct {
	Cause any
	Stack string
}

func (pe PanickedError) Error() string {
	return fmt.Sprintf("function panicked: %v", pe.Cause)
}

func RecoverIfPanic[T any](resultReference *Result[T]) {
	rec := recover()
	if rec != nil {

		stack := ""

		if CaptureStackOnPanic {
			stack = string(debug.Stack())
		}

		*resultReference = ResultFailed[T](PanickedError{
			Cause: rec,
			Stack: stack,
		})
	}
}
