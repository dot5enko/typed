package typed

type PanickedError struct {
	Cause any
}

func (PanickedError) Error() string {
	return "function panicked"
}

func RecoverIfPanic[T any](resultReference *Result[T]) {
	rec := recover()
	if rec != nil {
		*resultReference = ResultFailed[T](PanickedError{Cause: rec})
	}
}
