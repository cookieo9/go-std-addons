package xiter

import (
	"iter"
	"testing"
)

func filterTestCase[T any](name string, source []T, want []T, f func(T) bool) GenericTestCase {
	return IteratorCollectTest(name, Filter(sliceValues(source), f), want)
}

func excludeTestCase[T any](name string, source []T, want []T, f func(T) bool) GenericTestCase {
	return IteratorCollectTest(name, Exclude(sliceValues(source), f), want)
}

func TestFilterExclude(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	isEven := func(i int) bool { return i%2 == 0 }
	all := func(i int) bool { return true }

	t.Run("Filter", func(t *testing.T) {
		GenericTestCases{
			filterTestCase("all", nums, nums, all),
			filterTestCase("allEmpty", []int{}, []int{}, all),
			filterTestCase("allNil", nil, nil, all),

			filterTestCase("isEven", nums, []int{2, 4}, isEven),
			filterTestCase("isEvenEmpty", []int{}, []int{}, isEven),
			filterTestCase("isEvenNil", nil, nil, isEven),

			PanicTestCases(func(s iter.Seq[int]) iter.Seq[int] {
				return Filter(s, func(i int) bool { return true })
			}),
		}.Run(t)
	})

	t.Run("Exclude", func(t *testing.T) {
		GenericTestCases{
			excludeTestCase("all", nums, []int{}, all),
			excludeTestCase("allEmpty", []int{}, []int{}, all),
			excludeTestCase("allNil", nil, nil, all),

			excludeTestCase("isEven", nums, []int{1, 3, 5}, isEven),
			excludeTestCase("isEvenEmpty", []int{}, []int{}, isEven),
			excludeTestCase("isEvenNil", nil, nil, isEven),

			PanicTestCases(func(s iter.Seq[int]) iter.Seq[int] {
				return Exclude(s, func(i int) bool { return false })
			}),
		}.Run(t)
	})
}
