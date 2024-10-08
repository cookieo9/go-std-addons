package xiter

import "iter"

// One returns a function that yields the provided value t once.
func One[T any](t T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(t)
	}
}

// Forever returns a function that yields the provided value t indefinitely.
func Forever[T any](t T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(t) {
		}
	}
}

// Repeat returns a function that yields the provided value t n times.
func Repeat[T any](t T, n int) iter.Seq[T] {
	return Limit(Forever(t), n)
}
