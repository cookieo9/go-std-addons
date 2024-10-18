package xiter

import (
	"slices"
	"testing"
)

func TestUnique(t *testing.T) {
	piDigits := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	GenericTestCases{
		SliceCollectTest("empty", Unique(slices.Values([]int{})), nil),
		SliceCollectTest("pidigits", Unique(slices.Values(piDigits)), []int{3, 1, 4, 5, 9, 2, 6}),
		SliceCollectTest("pidigits-lim4", Limit(Unique(slices.Values(piDigits)), 4), []int{3, 1, 4, 5}),
		SliceCollectTest("all-unique", Unique(slices.Values([]int{1, 2, 3, 4, 5})), []int{1, 2, 3, 4, 5}),
	}.Run(t)
}
