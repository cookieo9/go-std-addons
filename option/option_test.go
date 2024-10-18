package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	x := Some(42)

	v, ok := x.Get()
	assert.Equal(t, v, 42)
	assert.Equal(t, x.Value(), 42)
	assert.Equal(t, x.GetOr(10), 42)
	assert.True(t, ok)
	assert.True(t, x.Ok())
	assert.NotPanics(t, func() {
		assert.Equal(t, x.Require(), 42)
	})

	y := None[string]()
	u, ok := y.Get()
	assert.Equal(t, u, "")
	assert.Equal(t, y.Value(), "")
	assert.Equal(t, y.GetOr("hello"), "hello")
	assert.False(t, ok)
	assert.False(t, y.Ok())
	assert.Panics(t, func() {
		t.Log("got file:", y.Require())
	})

	z := Map(x, func(v int) int { return v * 2 })
	v, ok = z.Get()
	assert.Equal(t, v, 84)
	assert.Equal(t, z.Value(), 84)
	assert.Equal(t, z.GetOr(10), 84)
	assert.True(t, ok)
	assert.True(t, z.Ok())
	assert.NotPanics(t, func() {
		assert.Equal(t, z.Require(), 84)
	})

	q := Map(y, func(v string) int { return len(v) + 5 })
	v, ok = q.Get()
	assert.Equal(t, v, 0)
	assert.Equal(t, q.Value(), 0)
	assert.Equal(t, q.GetOr(10), 10)
	assert.False(t, ok)
	assert.False(t, q.Ok())
	assert.Panics(t, func() {
		t.Log("got size:", q.Require())
	})

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
