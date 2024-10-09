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
		{Of(1, 2), Of(1, 2), true},
		{Of(1, 2), Of(1, 3), false},
		{Of(1, 2), Of(2, 2), false},
		{Of(1, 2), Of(2, 1), false},

		{Of(int8(1), int16(2)), Of(int8(1), int16(2)), true},
		{Of(int16(1), int8(2)), Of(int8(1), int16(2)), false},

		{Of[any, any](1, 2), Of(1, 2), false},
		{Of[any, any](1, 2), Of[any, any](1, 2), true},
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
		{Of("a", 1), Of("a", 2), -1},
		{Of("a", 1), Of("b", 1), -1},
		{Of("a", 1), Of("a", 1), 0},
		{Of("b", 1), Of("a", 1), 1},
		{Of("a", 2), Of("a", 1), 1},
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
	assert := assert.New(t)
	p0 := Of(1, "2")

	// Check accessor functions
	x, y := p0.Unpack()
	assert.Exactly(p0.A, p0.First(), "p0.A != p0.First()")
	assert.Exactly(p0.B, p0.Second(), "p0.B != p0.Second()")
	assert.Exactly(x, p0.A, "p0.A != p0.Unpack()[0]")
	assert.Exactly(y, p0.B, "p0.B != p0.Unpack()[1]")

	// Check that the fields are swapped
	p1 := p0.Swap()
	a, b := p0.Unpack()
	c, d := p1.Unpack()
	assert.Exactly(a, d, "p0.A != p1.B")
	assert.Exactly(b, c, "p0.B != p1.A")

	// Check that the String method matches format "%v"
	p0s := p0.String()
	p0f := fmt.Sprintf("%v", p0)
	assert.Exactly(p0s, p0f, "p0.String() != fmt.Sprintf(\"%v\", p0)")
}
