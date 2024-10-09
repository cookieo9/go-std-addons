// Package pair provides a simple Pair type to hold two values of potentially
// different types.
package pair

import (
	"fmt"
	"reflect"
)

// A Pair holds two values of potentially different types.
type Pair[T, U any] struct {
	A T
	B U
}

// String returns a string representation of the Pair in the format
// "Pair{<First>, <Second>}".
func (p Pair[T, U]) String() string {
	return fmt.Sprint(p)
}

// Format writes a string representation of the Pair in the format
// "Pair{<First>,<Second>}" as part of a fmt.*printf call. The format verb is
// ignored. If the # flag is set, the type parameters are also included in the
// format "Pair[T,U]{<First>,<Second>}".
func (p Pair[T, U]) Format(f fmt.State, _ rune) {
	a, b := p.Unpack()
	if f.Flag('#') {
		t := reflect.TypeOf(p)
		tT := t.Field(0).Type.String()
		uT := t.Field(1).Type.String()
		fmt.Fprintf(f, "Pair[%s,%s]{%#v,%#v}", tT, uT, a, b)
		return
	}
	fmt.Fprintf(f, "Pair{%#v,%#v}", a, b)
}

// First returns the first value stored in the Pair.
func (p Pair[T, U]) First() T {
	return p.A
}

// Second returns the second value stored in the Pair.
func (p Pair[T, U]) Second() U {
	return p.B
}

// Unpack returns the two values stored in the Pair.
func (p Pair[T, U]) Unpack() (T, U) {
	return p.A, p.B
}

// Swap returns a new Pair with the First and Second values swapped.
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{A: p.B, B: p.A}
}

// Of creates a new Pair with the given values.
func Of[T, U any](a T, b U) Pair[T, U] {
	return Pair[T, U]{A: a, B: b}
}

// Equal compares two Pair values for equality. It returns true if the two
// Pairs have the same First and Second values, and false otherwise.
//
// This only works when both fields implement the comparable interface.
func Equal[T, U comparable](a, b Pair[T, U]) bool {
	return a == b
}

// Less compares two Pair values by comparing their fields using the First
// field as the most significant field, and falling through to comparing the
// Second field if the First fields are equal.
//
// This only works when both fields implement the Ordered interface.
func Less[T, U Ordered](a, b Pair[T, U]) bool {
	return less2(a.A, b.A, a.B, b.B)
}

// Compare provides an integer representing the relative order of two Pair
// values. The semantics of the result are the same as for cmp.Compare, but
// the First fields are compared first, and only if they are equal are the
// Second fields compared.
//
// This only works when both fields implement the Ordered interface.
func Compare[T, U Ordered](a, b Pair[T, U]) int {
	return compare2(a.A, b.A, a.B, b.B)
}
