package xiter

import (
	"iter"
	"testing"
)

func mapTestCase[In, Out any](name string, source []In, want []Out, f func(In) Out) *SimpleTestCase[[]Out] {
	return IteratorCollectTest(name, Map(sliceValues(source), f), want)
}

func TestMap(t *testing.T) {
	doubleInt := func(i int) int { return i * 2 }
	halfFloat := func(f float64) float64 { return f / 2 }
	intToFloat := func(i int) float64 { return float64(i) }
	intToHalfFloat := func(i int) float64 { return float64(i) / 2 }

	GenericTestCases{
		mapTestCase("doubleInt", []int{1, 2, 3}, []int{2, 4, 6}, doubleInt),
		mapTestCase("doubleIntEmpty", []int{}, []int{}, doubleInt),
		mapTestCase("doubleIntNil", nil, nil, doubleInt),

		mapTestCase("halfFloat", []float64{1, 2, 3}, []float64{0.5, 1, 1.5}, halfFloat),
		mapTestCase("intToFloat", []int{1, 2, 3}, []float64{1, 2, 3}, intToFloat),
		mapTestCase("intToHalfFloat", []int{1, 2, 3}, []float64{0.5, 1, 1.5}, intToHalfFloat),

		PanicTestCases(func(s iter.Seq[int]) iter.Seq[int] {
			return Map(s, func(x int) int { return x })
		}),
	}.Run(t)
}
