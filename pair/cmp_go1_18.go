//go:build !go1.21
// +build !go1.21

package pair

// Ordered is a type constraint which represents any type that implments the
// standard comparison operators (<, >, <=, >=, ==, !=).
//
// This is a copy of the stdlib's cmp.Ordered type from go1.21+ for use in
// go1.18 - go1.20.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

func less[T Ordered](a, b T) bool {
	return ((a != a) && (b != b)) || a < b
}

func less2[T, U Ordered](a, b T, c, d U) bool {
	return less(a, b) || (a == b && less(c, d))
}

func compare[T Ordered](a, b T) int {
	if less(a, b) {
		return -1
	}
	if less(b, a) {
		return +1
	}
	return 0
}

func compare2[T, U Ordered](a, b T, c, d U) int {
	if x := compare(a, b); x != 0 {
		return x
	}
	return compare(c, d)
}
