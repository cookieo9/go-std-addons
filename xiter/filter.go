package xiter

// Filter applies the given predicate function f to each element in the input
// iterator it, and returns a new iterator that yields only the elements for
// which f returns true.
func Filter[T any](it func(func(T) bool), f func(T) bool) func(func(T) bool) {
	return process(it, func(t T, yield func(T) bool) bool {
		return !f(t) || yield(t)
	})
}

// Exclude applies the given predicate function f to each element in the input
// iterator it, and returns a new iterator that yields only the elements for
// which f returns false.
func Exclude[T any](it func(func(T) bool), f func(T) bool) func(func(T) bool) {
	return process(it, func(t T, yield func(T) bool) bool {
		return f(t) || yield(t)
	})
}
