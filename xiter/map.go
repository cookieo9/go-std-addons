package xiter

import "iter"

// Map applies the given function f to each element of the input iterator, and
// yields the result of each to the output iterator.
func Map[T, U any](it iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for t := range it {
			if !yield(f(t)) {
				return
			}
		}
	}
}

// MapOut applies the given function f to each element of the input iterator,
// producing two values that are yielded by the output iterator.
func MapOut[T, V, W any](it iter.Seq[T], f func(T) (V, W)) iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {
		for t := range it {
			if !yield(f(t)) {
				return
			}
		}
	}
}

// MapIn applies the given function f to each pair of elements from the input
// iterator, producing a single value that is yielded by the output iterator.
func MapIn[T, U, V any](it iter.Seq2[T, U], f func(T, U) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for t, u := range it {
			if !yield(f(t, u)) {
				return
			}
		}
	}
}
