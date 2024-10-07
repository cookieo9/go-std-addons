package xiter

import (
	"iter"
	"slices"
	"testing"
)

func TestOne(t *testing.T) {
	GenericTestCases{
		SliceCollectTest("One(1)", One(1), []int{1}),
		SliceCollectTest("One(`five`)", One("five"), []string{"five"}),

		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return One(42)
		}),
	}.Run(t)
}

func TestRepeat(t *testing.T) {
	GenericTestCases{
		SliceCollectTest("Repeat(1,0)", Repeat(1, 0), []int{}),
		SliceCollectTest("Repeat(1,1)", Repeat(1, 1), []int{1}),
		SliceCollectTest("Repeat(1,2)", Repeat(1, 2), []int{1, 1}),

		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Repeat(42, 1)
		}),
	}.Run(t)
}

func TestForever(t *testing.T) {
	GenericTestCases{
		SliceCollectTest("Limit(Forever(5.5),3)", Limit(Forever(5.5), 3), []float64{5.5, 5.5, 5.5}),

		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Limit(Forever(12), 42)
		}),
	}.Run(t)
}
