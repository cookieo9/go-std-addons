package xiter

import "iter"

func process[T, U any](it iter.Seq[T], f func(T, func(U) bool) bool) iter.Seq[U] {
	return func(yield func(U) bool) {
		it(func(t T) bool { return f(t, yield) })
	}
}

func one[T any](t T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(t)
	}
}

func sliceValues[T any](values []T) iter.Seq[T] {
	return process(one(values), func(t []T, yield func(T) bool) bool {
		for _, v := range t {
			if !yield(v) {
				return false
			}
		}
		return true
	})
}

func sliceCollect[T any](it iter.Seq[T]) []T {
	var out []T
	it(func(t T) bool {
		out = append(out, t)
		return true
	})
	return out
}
