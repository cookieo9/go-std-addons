package xiter

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func limitTestCase[T any](name string, src []T, n int, want []T) *SimpleTestCase[[]T] {
	return SliceCollectTest(name, Limit(slices.Values(src), n), want)
}

func TestLimit(t *testing.T) {
	GenericTestCases{
		limitTestCase("emptyNone", []int{}, 0, []int{}),
		limitTestCase("emptyOne", []int{}, 1, []int{}),
		limitTestCase("emptyMany", []int{}, 42, []int{}),

		limitTestCase("nilNone", nil, 0, []int{}),
		limitTestCase("nilOne", nil, 1, []int{}),
		limitTestCase("nilMany", nil, 42, []int{}),

		limitTestCase("someNone", []int{1, 2, 3}, 0, []int{}),
		limitTestCase("someOne", []int{1, 2, 3}, 1, []int{1}),
		limitTestCase("someTwo", []int{1, 2, 3}, 2, []int{1, 2}),
		limitTestCase("someMany", []int{1, 2, 3}, 42, []int{1, 2, 3}),

		PanicTestCases(func(s iter.Seq[bool]) iter.Seq[bool] {
			return Limit(s, 42)
		}),
	}.Run(t)
}

func whileTestCase[T any](name string, src []T, f func(T) bool, want []T) *SimpleTestCase[[]T] {
	return SliceCollectTest(name, While(slices.Values(src), f), want)
}

func untilTestCase[T any](name string, src []T, f func(T) bool, want []T) *SimpleTestCase[[]T] {
	return SliceCollectTest(name, Until(slices.Values(src), f), want)
}

func TestWhileUntil(t *testing.T) {
	lessThan5 := func(i int) bool { return i < 5 }
	moreThan3 := func(i int) bool { return i > 3 }
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	GenericTestCases{
		whileTestCase("emptyWhile", []int{}, lessThan5, []int{}),
		whileTestCase("nilWhile", nil, lessThan5, []int{}),
		untilTestCase("emptyUntil", []int{}, lessThan5, []int{}),
		untilTestCase("nilUntil", nil, lessThan5, []int{}),

		whileTestCase("whileLessThan5", nums, lessThan5, []int{1, 2, 3, 4}),
		untilTestCase("untilLessThan5", nums, lessThan5, []int{}),

		whileTestCase("whileMoreThan3", nums, moreThan3, []int{}),
		untilTestCase("untilMoreThan3", nums, moreThan3, []int{1, 2, 3}),

		PanicTestCases(func(f iter.Seq[int]) iter.Seq[int] {
			return While(f, func(i int) bool { return i > 10 })
		}),
		PanicTestCases(func(f iter.Seq[int]) iter.Seq[int] {
			return Until(f, func(i int) bool { return i > 10 })
		}),
	}.Run(t)
}

type flResult[T any] struct {
	Value T
	Ok    bool
}

func makeFirstTest[T any](name string, src []T, want T, ok bool) *SimpleTestCase[flResult[T]] {
	tc := SimpleTest(name, func(t *testing.T) flResult[T] {
		v, ok := First(slices.Values(src))
		return flResult[T]{v, ok}
	})
	return tc.Compare(flResult[T]{want, ok}, assert.Equal).Args("match result")
}

func makeLastTest[T any](name string, src []T, want T, ok bool) *SimpleTestCase[flResult[T]] {
	tc := SimpleTest(name, func(t *testing.T) flResult[T] {
		v, ok := Last(slices.Values(src))
		return flResult[T]{v, ok}
	})
	return tc.Compare(flResult[T]{want, ok}, assert.Equal).Args("match result")
}

func TestFirstLast(t *testing.T) {
	GenericTestCases{
		makeFirstTest("emptyFirst", []int{}, 0, false),
		makeFirstTest("nilFirst", nil, 0, false),
		makeLastTest("emptyLast", []int{}, 0, false),
		makeLastTest("nilLast", nil, 0, false),

		makeFirstTest("someFirst", []int{1, 2, 3}, 1, true),
		makeLastTest("someLast", []int{1, 2, 3}, 3, true),

		PanicTestCases(func(f iter.Seq[int]) iter.Seq[int] {
			First(f)
			return f
		}),
		PanicTestCases(func(f iter.Seq[int]) iter.Seq[int] {
			Last(f)
			return f
		}),
	}.Run(t)
}
