package xiter

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterializeCount(t *testing.T) {
	source := slices.Values([]int{1, 2, 3})
	source, n := CountUses(source)
	m := Materialize(source)
	assert.Equal(t, 0, *n, "source iterator not yet iterated")
	slices.Collect(m)
	assert.Equal(t, 1, *n, "materialized iterator causes source to be iterated")
	for range 50 {
		slices.Collect(m)
		assert.Equal(t, 1, *n, "materialized iterator prevents further iteration of source")
	}
}

func TestMaterializeOnce(t *testing.T) {
	data := []int{1, 2, 3}
	want := slices.Clone(data)
	source := func(yield func(int) bool) {
		var x int
		for len(data) > 0 {
			x, data = data[0], data[1:]
			if !yield(x) {
				data = nil
				return
			}
		}
	}
	m := Materialize(source)
	for range 50 {
		got := slices.Collect(m)
		assert.Equal(t, want, got, "same sequence")
	}
	gotEmpty := slices.Collect(source)
	assert.Empty(t, gotEmpty, "source iterator is exhausted")
}

func TestMaterializePanic(t *testing.T) {
	PanicTestCases[int](Materialize[int]).Run(t)
}

func TestMaterializeShort(t *testing.T) {
	values := []int{1, 2, 3}
	m := Materialize(slices.Values(values))
	mLimited := Limit(m, 1)

	for range 10 {
		wantLimited := []int{1}
		got := slices.Collect(mLimited)
		assert.Equal(t, wantLimited, got, "got limited sequence")

		wantFull := slices.Clone(values)
		got = slices.Collect(m)
		assert.Equal(t, wantFull, got, "got full original sequence")
	}
}
