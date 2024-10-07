package xiter

import "iter"

func process[T, U any](it iter.Seq[T], f func(T, func(U) bool) bool) iter.Seq[U] {
	return func(yield func(U) bool) {
		it(func(t T) bool { return f(t, yield) })
	}
}
