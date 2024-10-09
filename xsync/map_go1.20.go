//go:build go1.20
// +build go1.20

package xsync

// CompareAndDelete compares the value associated with the given key to the provided old value,
// and if they are equal, deletes the key-value pair from the map. It returns whether the
// key-value pair was deleted.
//
// Requires Go 1.20+.
func (m *Map[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	return m.m.CompareAndDelete(k, old)
}

// CompareAndSwap compares the value associated with the given key to the provided old value,
// and if they are equal, sets the value to the new value. It returns whether the value was swapped.
//
// Requires Go 1.20+.
func (m *Map[K, V]) CompareAndSwap(k K, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(k, old, new)
}

// Swap sets the value associated with the given key to the new value
// returning the previous value (if any) and true if the swap replaced an
// existing value.
//
// Requires Go 1.20+.
func (m *Map[K, V]) Swap(k K, v V) (previous V, swapped bool) {
	if p, swapped := m.m.Swap(k, v); swapped {
		return p.(V), true
	}
	return *new(V), false
}
