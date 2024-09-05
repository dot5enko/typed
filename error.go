package typed

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

type PanicError struct {
	message string
	stack   []string
}

func (pe *PanicError) Error() string {
	return pe.message
}

func (pe *PanicError) Stack() []string {
	return pe.stack
}

func (pe *PanicError) FormatStack() string {

	jb, _ := json.MarshalIndent(pe.stack, " ", " ")
	return string(jb)
}

type HasStack interface {
	Stack() []string
}

var DefaultStackCaptureDepth = 10

func NewPanicErrorConfig(data any, lines, skip int) *PanicError {
	var msg string

	switch converted := data.(type) {
	case string:
		msg = converted
	default:
		msg = fmt.Sprintf("%v", data)
	}

	return &PanicError{
		message: msg,
		stack:   RecordDebugStackConfig(lines, skip),
	}
}

func NewPanicError(data any) *PanicError {
	return NewPanicErrorConfig(data, DefaultStackCaptureDepth, 2) // skip this hop + caller itself
}

func RecordDebugStack() []string {
	return RecordDebugStackConfig(DefaultStackCaptureDepth, 1) // skip current hop
}

// add skip lines arg
func RecordDebugStackConfig(callers, skip int) []string {
	stack := []string{}

	if callers <= 0 {
		callers = DefaultStackCaptureDepth
	}

	if skip <= 0 {
		skip = 0
	}

	collected := 0

	i := 0

	for {

		if collected >= callers {
			break
		}

		i++

		curi := i - 1

		if curi < skip {
			continue
		}

		_, file, line, ok := runtime.Caller(i)
		if ok {
			stack = append(stack, fmt.Sprintf(" %s : %d", file, line))
			collected += 1
		} else {
			break
		}
	}

	return stack
}

// records stack trace
func RecoverPanic(onPanic func(pe *PanicError)) {
	rec := recover()

	if rec != nil {
		onPanic(NewPanicErrorConfig(rec, DefaultStackCaptureDepth, 3)) // defer, recover panic, self
	}
}

func RecoverPanicToLog() {
	RecoverPanic(func(pe *PanicError) {
		log.Printf("recovered  : %s", pe.Error())
		pe.FormatStack()
	})
}

func GetStackForError(e error) []string {

	switch convertedErr := e.(type) {
	case HasStack:
		return convertedErr.Stack()
	default:
		return RecordDebugStackConfig(DefaultStackCaptureDepth, 2) // skip caller and current
	}

}

func SafeGoroutine(cb func()) {
	SafeGoroutineWithCb(cb, nil)
}

func WrapPanic(cb func()) func() {
	return func() {
		defer RecoverPanic(func(pe *PanicError) {
			log.Printf("safe goroutine paniced: %s", pe.Error())
			pe.FormatStack()
		})

		cb()
	}
}

func SafeGoroutineWithCb(cb func(), onPanic func(pe *PanicError)) {
	go func() {
		defer RecoverPanic(func(pe *PanicError) {
			if onPanic != nil {
				onPanic(pe)
			} else {
				log.Printf("safe goroutine paniced: %s", pe.Error())
				pe.FormatStack()
			}
		})

		cb()
	}()
}

func Isolate(task func()) (err error) {

	defer RecoverPanic(func(pe *PanicError) {
		err = pe
	})

	task()

	return nil

}
