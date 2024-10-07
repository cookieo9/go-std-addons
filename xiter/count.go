package xiter

// Countable represents the types usable by Count* and Range* functions. They
// are the types that can be incremented or decremented, as well as compared
// for greater than or less than.
type Countable interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Count returns an iterator that yields successive values starting from the
// given start value, incrementing by 1 each time. It will continue forever
// and wrap on overflow. It is equivalent to calling CountUp with a step of 1.
func Count[T Countable](start T) func(func(T) bool) {
	return CountUp(start, 1)
}

// CountUp returns an iterator that yields successive values starting from the
// given start value, incrementing by the given step each time. It will continue
// forever and wrap on overflow.
func CountUp[T Countable](start, step T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for i := start; yield(i); i += step {
		}
	}
}

// CountDown returns an iterator that yields successive values starting from the
// given start value, decrementing by the given step each time. It will continue
// forever and wrap on underflow.
func CountDown[T Countable](start, step T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for i := start; yield(i); i -= step {
		}
	}
}

// Range returns an iterator that yields successive values from start to finish
// (exclusive) with a step of 1. It is equivalent to calling RangeBy with a
// step of 1.
func Range[T Countable](start, end T) func(func(T) bool) {
	return RangeBy(start, end, 1)
}

// RangeBy returns an iterator that yields successive values from start to end
// with the given step size. If start is less than end, the iterator counts up,
// otherwise it counts down. When counting up, the iterator will stop when
// i >= end, and when counting down, it will stop when i <= end.
func RangeBy[T Countable](start, end, step T) func(func(T) bool) {
	if start < end {
		return While(CountUp(start, step), func(i T) bool { return i < end })
	}
	return While(CountDown(start, step), func(i T) bool { return i > end })
}
