// Package xsync provides generic wrapper to build-in concurrent types in the
// sync package of the standard library.
package xsync

import "sync"

// Map provides a type-safe wrapper around sync.Map. As a result, it is safe to
// use concurrently and the zero value is ready to use.
//
// Aside from the types, all existing APIs are the same as sync.Map for your
// version of Go. i.e.: m.Clear() will only exist if the Go version is 1.23+.
//
// See https://pkg.go.dev/sync#Map for more details on the underlying container.
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load retrieves the value stored in the map for a key, or returns the zero
// value of the value type and false if no value is present.
func (m *Map[K, V]) Load(k K) (value V, loaded bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return *new(V), false
	}
	return v.(V), true
}

// Store sets the value for a key in the map.
func (m *Map[K, V]) Store(k K, v V) {
	m.m.Store(k, v)
}

// LoadOrStore loads the value stored in the map for a key, or stores and
// returns the given value if no value is present. If the value was loaded,
// then the second return value is true, otherwise on store it is false.
func (m *Map[K, V]) LoadOrStore(k K, v V) (value V, loaded bool) {
	v2, load := m.m.LoadOrStore(k, v)
	return v2.(V), load
}

// Delete removes the value for a key from the map.
func (m *Map[K, V]) Delete(k K) {
	m.m.Delete(k)
}

// Range calls the provided function f for each key and value present in the map.
// The function f is called for each element as long as f returns true.
// The iteration order is not specified and is not guaranteed to be the same
// across multiple iterations.
func (m *Map[K, V]) Range(f func(k K, v V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// LoadAndDelete atomically loads and deletes the value stored in the map for a
// key. It returns the loaded value and a boolean indicating whether the key was
// present. If the key was not present, the returned value will be the zero
// value for the value type.
func (m *Map[K, V]) LoadAndDelete(k K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(k)
	if !ok {
		return *new(V), false
	}
	return v.(V), true
}
