//go:build go1.23
// +build go1.23

package xsync

// Clear removes all key-value pairs from the Map.
//
// Requires Go 1.23+.
func (m *Map[K, V]) Clear() {
	m.m.Clear()
}
