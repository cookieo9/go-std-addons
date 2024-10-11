// Package xerrors provides utilities for working with code that may panic
// with an error value.
package xerrors

// Catch is a utility function that calls the provided function f and recovers
// any panics that occur should they contain an error. If an error is recovered,
// it will be returned. If the panic doesn't contain an error, it will be
// re-thrown. If no panic occurs, nil will be returned. This is useful for
// simplifying code that needs to call functions that may panic with errors,
// such as running a range loop over an iterator that panics on error.
func Catch(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()
	f()
	return nil
}

// CatchValue is a utility function that calls the provided function f and returns
// the value it returns or an error if the function panics with an error. If the
// panic doesn't contain an error, it will be re-thrown. If no panic occurs, then
// the return value of the function will be returned with a nil error.
//
// This is useful for calling functions that may panic while returning a value,
// such as functions that process an iterator and return a value.
func CatchValue[T any](f func() T) (value T, err error) {
	err = Catch(func() { value = f() })
	return value, err
}

// Must is a utility function that takes a value and an error, and returns the value.
// If the error is not nil, Must will panic with the error.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
