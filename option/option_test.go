package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	x := Some(42)
	assertSome(t, x, 42, 101)

	y := None[string]()
	assertNone(t, y, "hello")

	z := Map(x, func(v int) int { return v * 2 })
	assertSome(t, z, 84, 102)

	q := Map(y, func(v string) int { return len(v) + 5 })
	assertNone(t, q, 103)

	done := false
	y.Do(func(v string) {
		done = true
	})
	assert.False(t, done)

	z.Do(func(v int) {
		if v > 50 {
			done = true
		}
	})
	assert.True(t, done)
}

// assertSome checks that the given Value contains the expected value, and that
// the various accessor methods (Get, GetValue, GetOr, Ok, Require) behave as
// expected for a Some value.
func assertSome[T any](t *testing.T, o Value[T], expectedValue, defaultValue T) {
	v, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, v, expectedValue)
	assert.Equal(t, o.GetValue(), expectedValue)
	assert.Equal(t, o.GetOr(defaultValue), expectedValue)
	assert.True(t, o.Ok())
	assert.NotPanics(t, func() {
		assert.Equal(t, o.Require(), expectedValue)
	})
}

// assertNone checks that the given Value contains no value, and that the various
// accessor methods (Get, GetValue, GetOr, Ok, Require) behave as expected for a
// None value.
func assertNone[T any](t *testing.T, o Value[T], defaultValue T) {
	var zeroValue T
	v, ok := o.Get()
	assert.False(t, ok)
	assert.Equal(t, v, zeroValue)
	assert.Equal(t, o.GetValue(), zeroValue)
	assert.Equal(t, o.GetOr(defaultValue), defaultValue)
	assert.False(t, o.Ok())
	assert.Panics(t, func() {
		assert.Equal(t, o.Require(), zeroValue)
	})
}
