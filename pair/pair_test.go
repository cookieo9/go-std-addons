package pair

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquality(t *testing.T) {
	testCases := []struct {
		a, b  any
		equal bool
	}{
		{New(1, 2), New(1, 2), true},
		{New(1, 2), New(1, 3), false},
		{New(1, 2), New(2, 2), false},
		{New(1, 2), New(2, 1), false},

		{New(int8(1), int16(2)), New(int8(1), int16(2)), true},
		{New(int16(1), int8(2)), New(int8(1), int16(2)), false},

		{New[any, any](1, 2), New(1, 2), false},
		{New[any, any](1, 2), New[any, any](1, 2), true},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.equal {
			name := fmt.Sprintf("%#v == %#v", tc.a, tc.b)
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, tc.a, tc.b)
			})
		} else {
			name := fmt.Sprintf("%#v != %#v", tc.a, tc.b)
			t.Run(name, func(t *testing.T) {
				assert.NotEqual(t, tc.a, tc.b)
			})
		}
	}
}

func TestOrdering(t *testing.T) {
	testCases := []struct {
		a, b Pair[string, int]
		want int
	}{
		{New("a", 1), New("a", 2), -1},
		{New("a", 1), New("b", 1), -1},
		{New("a", 1), New("a", 1), 0},
		{New("b", 1), New("a", 1), 1},
		{New("a", 2), New("a", 1), 1},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Compare(%v,%v)", tc.a, tc.b), func(t *testing.T) {
			assert.Equal(t, tc.want, Compare(tc.a, tc.b))
			if tc.want < 0 {
				assert.True(t, Less(tc.a, tc.b))
			} else if tc.want > 0 {
				assert.False(t, Less(tc.a, tc.b))
			} else {
				assert.False(t, Less(tc.a, tc.b))
				assert.False(t, Less(tc.b, tc.a))
				assert.True(t, Equal(tc.a, tc.b))
				assert.True(t, Equal(tc.b, tc.a))
			}
		})
	}
}

func TestInvariants(t *testing.T) {
	p0 := New(1, "2")

	// Check accessor functions
	x, y := Unpack(p0)
	assert.Exactly(t, p0.First, First(p0), "p0.First != First(p0)")
	assert.Exactly(t, p0.Second, Second(p0), "p0.Second != Second(p0)")
	assert.Exactly(t, x, p0.First, "p0.First != Unpack(p0)[0]")
	assert.Exactly(t, y, p0.Second, "p0.Second != Unpack(p0)[1]")

	// Check that the fields are swapped
	p1 := Swap(p0)
	a, b := Unpack(p0)
	c, d := Unpack(p1)
	assert.Exactly(t, a, d, "p0.First != p1.Second")
	assert.Exactly(t, b, c, "p0.Second != p1.First")

	// Check that the String method matches format "%v"
	p0s := p0.String()
	p0f := fmt.Sprintf("%v", p0)
	assert.Exactly(t, p0s, p0f, "p0.String() != fmt.Sprintf(\"%v\", p0)")
}
