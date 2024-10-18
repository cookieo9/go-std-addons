package xiter

import "iter"

// Unique returns a new sequence that contains only the unique elements from the
// input sequence. The elements must be comparable.
func Unique[T comparable](in iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for v := range in {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}
