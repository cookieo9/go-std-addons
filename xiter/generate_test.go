package xiter

import "testing"

func TestOne(t *testing.T) {
	GenericTestCases{
		IteratorCollectTest("One(1)", One(1), []int{1}),
		IteratorCollectTest("One(`five`)", One("five"), []string{"five"}),

		PanicTestCases(func(it func(func(int) bool)) func(func(int) bool) {
			sliceCollect(it)
			return One(42)
		}),
	}.Run(t)
}

func TestRepeat(t *testing.T) {
	GenericTestCases{
		IteratorCollectTest("Repeat(1,0)", Repeat(1, 0), []int{}),
		IteratorCollectTest("Repeat(1,1)", Repeat(1, 1), []int{1}),
		IteratorCollectTest("Repeat(1,2)", Repeat(1, 2), []int{1, 1}),

		PanicTestCases(func(it func(func(int) bool)) func(func(int) bool) {
			sliceCollect(it)
			return Repeat(42, 1)
		}),
	}.Run(t)
}

func TestForever(t *testing.T) {
	GenericTestCases{
		IteratorCollectTest("Limit(Forever(5.5),3)", Limit(Forever(5.5), 3), []float64{5.5, 5.5, 5.5}),

		PanicTestCases(func(it func(func(int) bool)) func(func(int) bool) {
			sliceCollect(it)
			return Limit(Forever(12), 42)
		}),
	}.Run(t)
}
