package xiter

// Map applies the given function f to each element of the input iterator it,
// and returns a new iterator that yields the results of applying f.
func Map[T, U any](it func(func(T) bool), f func(T) U) func(func(U) bool) {
	return process(it, func(t T, yield func(U) bool) bool {
		return yield(f(t))
	})
}
