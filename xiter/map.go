package xiter

import "iter"

// Map applies the given function f to each element of the input iterator it,
// and returns a new iterator that yields the results of applying f.
func Map[T, U any](it iter.Seq[T], f func(T) U) iter.Seq[U] {
	return Process(it, func(t T, yield func(U) bool) bool {
		return yield(f(t))
	})
}
