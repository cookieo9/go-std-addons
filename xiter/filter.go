package xiter

import "iter"

// Filter applies the given predicate function f to each element in the input
// iterator it, and returns a new iterator that yields only the elements for
// which f returns true.
func Filter[T any](it iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return process(it, func(t T, yield func(T) bool) bool {
		return !f(t) || yield(t)
	})
}

// Exclude applies the given predicate function f to each element in the input
// iterator it, and returns a new iterator that yields only the elements for
// which f returns false.
func Exclude[T any](it iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return process(it, func(t T, yield func(T) bool) bool {
		return f(t) || yield(t)
	})
}
