package typed

import (
	"fmt"
	"reflect"
)

type Result[T any] struct {
	has bool
	val *T
	err error
}

func (r Result[T]) String() string {

	obj := new(T)

	return fmt.Sprintf("#result[%s]{has: %v}", reflect.ValueOf(obj).Type().Name(), r.has)
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	obj := new(T)
	return []byte(fmt.Sprintf("\"#result[%s]{has: %v}\"", reflect.ValueOf(obj).Type().Name(), r.has)), nil
}

func (opt Result[T]) IsOk() bool {
	return opt.has
}

func (opt Result[T]) Unwrap() T {

	if !opt.has {
		panic("unwrapping failed result")
	}

	return *opt.val
}

func (opt Result[T]) UnwrapDefault() T {

	if opt.IsOk() {
		return opt.Unwrap()
	} else {
		var def T
		return def
	}
}

func (opt Result[T]) UnwrapOrDefault(def T) T {

	if opt.IsOk() {
		return opt.Unwrap()
	} else {
		return def
	}
}

func (opt Result[T]) UnwrapOrPanic(msg string) T {
	if opt.IsOk() {
		return opt.Unwrap()
	} else {
		panic(msg)
	}
}

func (opt Result[T]) UnwrapError() error {
	return opt.err
}

func (opt Result[T]) UnwrapClassic() (obj T, e error) {

	if opt.IsOk() {
		obj = opt.Unwrap()
	} else {
		e = opt.UnwrapError()
	}

	return
}

func (opt Result[T]) Accept(f func(*T)) (nextResult Result[T]) {

	defer func() {

		rec := recover()

		if rec != nil {
			nextResult = ResultFailed[T](fmt.Errorf("panicked while accepting result ok: %v", rec))
		}
	}()

	if opt.IsOk() {
		f(opt.val)
	}
	return opt
}

func (opt Result[T]) Then(f func(*T) *Result[T]) Result[T] {
	if opt.IsOk() {
		result := f(opt.val)
		if result != nil {
			return *result
		}
	}
	return opt
}

func (opt Result[T]) Fail(f func(e error)) Result[T] {
	if !opt.IsOk() {
		f(opt.err)
	}
	return opt
}

func (opt *Result[T]) SetOk(v T) *Result[T] {
	opt.val = &v
	opt.has = true
	return opt
}

func (opt *Result[T]) SetFail(v error) *Result[T] {
	opt.err = v
	opt.has = false
	return opt
}

func ResultOk[T any](v T) Result[T] {
	res := Result[T]{}
	res.SetOk(v)

	return res
}

func ResultFailed[T any](err error) Result[T] {
	res := Result[T]{}
	res.SetFail(err)
	return res
}
