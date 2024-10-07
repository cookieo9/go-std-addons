package xiter

import "sync"

// Materialize returns an iterator that yield the same values as the original
// iterator, using a cached copy of the data from the original iterator. The
// cache is generated lazily, the first time the iterator is iterated over, and
// is reused for all subsequent iterations. The input iterator will only be
// iterated over once by the Materialize iterator.
//
// The resulting iterator can be used multiple times, and can be used by
// parallel goroutines, even when the original iterator cannot.
//
// Warning: Do not use Materialize on an indefinite iterator, as the cache will
// grow indefinitely and consume all available memory.
func Materialize[T any](it func(func(T) bool)) func(func(T) bool) {
	values := sync.OnceValue(func() []T { return sliceCollect(it) })

	return func(yield func(T) bool) {
		sliceValues(values())(yield)
	}
}
