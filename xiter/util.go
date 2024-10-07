package xiter

func process[T, U any](it func(func(T) bool), f func(T, func(U) bool) bool) func(func(U) bool) {
	return func(yield func(U) bool) {
		it(func(t T) bool { return f(t, yield) })
	}
}

func one[T any](t T) func(func(T) bool) {
	return func(yield func(T) bool) {
		yield(t)
	}
}

func sliceValues[T any](values []T) func(func(T) bool) {
	return process(one(values), func(t []T, yield func(T) bool) bool {
		for _, v := range t {
			if !yield(v) {
				return false
			}
		}
		return true
	})
}

func sliceCollect[T any](it func(func(T) bool)) []T {
	var out []T
	it(func(t T) bool {
		out = append(out, t)
		return true
	})
	return out
}
