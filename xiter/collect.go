package xiter

import "iter"

// Collect is a higher-order function that applies a function f to update an
// accumulated value using each element of an iterator. It starts with an
// initial value of the accumulator and returns the final accumulated value.
// Using an iterator with no elements will return the initial accumulator value.
func Collect[T, Accum any](it iter.Seq[T], start Accum, f func(Accum, T) Accum) Accum {
	out := start
	for t := range it {
		out = f(out, t)
	}
	return out
}
