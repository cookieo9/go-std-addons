package xiter

import "iter"

// Process is a helper for generating new iterators from existing iterators. It
// accepts a function that will be called for each element in the iterator, and
// it is expected to use a provided funtion to yield any number of elements, or
// even none at all for the current element.
func Process[T, U any](it iter.Seq[T], f func(T, func(U) bool) bool) iter.Seq[U] {
	return func(yield func(U) bool) {
		it(func(t T) bool { return f(t, yield) })
	}
}
