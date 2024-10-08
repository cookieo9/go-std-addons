//go:build go1.21
// +build go1.21

package pair

import (
	"cmp"
)

// Ordered is an alias for the cmp.Ordered type, which represents any type that
// implements the standard comparison operators (<, >, <=, >=, ==, !=).
type Ordered = cmp.Ordered

func less2[T, U Ordered](a, b T, c, d U) bool {
	return cmp.Less(a, b) || (a == b && cmp.Less(c, d))
}

func compare2[T, U Ordered](a, b T, c, d U) int {
	if x := cmp.Compare(a, b); x != 0 {
		return x
	}
	return cmp.Compare(c, d)
}
