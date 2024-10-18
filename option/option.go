package option

// Value is a typed wrapper for a value that may or may not be present. It is
// an immutable value type.
type Value[T any] struct {
	value T
	ok    bool
}

// Of returns a Value[T] that contains the given value v if ok is true, or a
// Value[T] that does not contain a value if ok is false. It will not store
// the value if ok is false.
func Of[T any](v T, ok bool) Value[T] {
	if ok {
		return Value[T]{value: v, ok: true}
	}
	return Value[T]{}
}

// Some returns a typed value that is present using the given value.
func Some[T any](x T) Value[T] {
	return Of(x, true)
}

// None returns a typed value that is not present.
func None[T any]() Value[T] {
	return Value[T]{}
}

// Get returns the value, and a boolean indicating whether the value is present.
// If the value is not present, the zero value of the type is returned.
func (o Value[T]) Get() (T, bool) {
	return o.value, o.ok
}

// Ok returns a boolean indicating whether the Value contains a valid value.
func (o Value[T]) Ok() bool {
	return o.ok
}

// Do applies the given function f to the value in the Value[T] if present.
// If the Value[T] is not present, Do does nothing.
func (o Value[T]) Do(f func(v T)) {
	if o.ok {
		f(o.value)
	}
}

// GetOr returns the value in the Value[T] if present, otherwise returns the
// provided default value v.
func (o Value[T]) GetOr(v T) T {
	if o.ok {
		return o.value
	}
	return v
}

// Value returns the value contained in the Value[T]. If the Value[T] is not
// present, the zero value of the type T is returned.
func (o Value[T]) Value() T {
	return o.value
}

// Require returns the value contained in the Value[T]. If the Value[T] is not
// present, it panics with the message "value is not present".
func (o Value[T]) Require() T {
	if !o.ok {
		panic("value is not present")
	}
	return o.value
}

// Map applies the given function f to the value in the Value[T] if present,
// and returns a new Value[U] containing the result. If the Value[T] is not
// present, Map returns a Value[U] that is not present.
func Map[T any, U any](v Value[T], f func(T) U) Value[U] {
	if v.ok {
		return Some(f(v.value))
	}
	return None[U]()
}
