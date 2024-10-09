package xsync

// Cache retrieves the value stored in the Map for the given key, or computes
// and stores the value using the provided function if no value is present. If
// the value was loaded, then the second return value is true, otherwise on
// compute/store it is false.
func Cache[K comparable, V any](m *Map[K, V], k K, g func() V) (value V, loaded bool) {
	if v, ok := m.Load(k); ok {
		return v, true
	}
	return m.LoadOrStore(k, g())
}

// MapLen computes and returns the number of elements in the given Map.
//
// Notes:
//   - This is not a constant-time operation. (Implemented via Range)
//   - This is not a safe value to rely on if the Map is being modified
//     concurrently, as the calculation doesn't prevent the size from changing
//     even while it is being calculated.
func MapLen[K comparable, V any](m *Map[K, V]) int {
	var n int
	m.m.Range(func(any, any) bool {
		n++
		return true
	})
	return n
}

// MapClear clears the contents of the given Map. If the Map implements a Clear
// method (i.e. in go1.23+), that is used. Otherwise, the Map is cleared by
// iterating over all keys and deleting them one by one.
func MapClear[K comparable, V any](m *Map[K, V]) {
	if clr, ok := any(m).(interface{ Clear() }); ok {
		clr.Clear()
		return
	}
	m.Range(func(k K, _ V) bool {
		m.Delete(k)
		return true
	})
}
