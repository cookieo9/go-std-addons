package xiter

// Limit returns a new iterator that yields at most n elements from the input
// iterator.
func Limit[T any](it func(func(T) bool), n int) func(func(T) bool) {
	return func(yield func(T) bool) {
		i := 0
		it(func(t T) bool {
			if i >= n {
				return false
			}
			i++
			return yield(t)
		})
	}
}

// While returns a new iterator function that yields elements from the given
// iterator function it as long as the provided condition function f returns
// true for each element.
func While[T any](it func(func(T) bool), f func(T) bool) func(func(T) bool) {
	return process(it, func(t T, yield func(T) bool) bool { return f(t) && yield(t) })
}

// Until returns a new iterator function that yields elements from the given
// iterator function it as long as the provided condition function f returns
// false for each element.
func Until[T any](it func(func(T) bool), f func(T) bool) func(func(T) bool) {
	return process(it, func(t T, yield func(T) bool) bool { return !f(t) && yield(t) })
}

// Last returns the last element yielded by the given iterator function,
// and a boolean value showing if there was an element. An empty iterator will
// result in a zero value and false.
func Last[T any](it func(func(T) bool)) (value T, ok bool) {
	it(func(t T) bool {
		value, ok = t, true
		return true
	})
	return value, ok
}

// First returns the first element yielded by the given iterator function,
// and a boolean value showing if there was an element. An empty iterator will
// result in a zero value and false.
func First[T any](it func(func(T) bool)) (value T, ok bool) {
	it(func(t T) bool {
		value, ok = t, true
		return false
	})
	return value, ok
}
