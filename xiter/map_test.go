package xiter

import (
	"iter"
	"math"
	"slices"
	"testing"

	"github.com/cookieo9/go-std-addons/pair"
)

func mapTestCase[In, Out any](name string, source []In, want []Out, f func(In) Out) *SimpleTestCase[[]Out] {
	return SliceCollectTest(name, Map(slices.Values(source), f), want)
}

func TestMap(t *testing.T) {
	doubleInt := func(i int) int { return i * 2 }
	halfFloat := func(f float64) float64 { return f / 2 }
	intToFloat := func(i int) float64 { return float64(i) }
	intToHalfFloat := func(i int) float64 { return float64(i) / 2 }

	TestSuite{
		mapTestCase("doubleInt", list(1, 2, 3), list(2, 4, 6), doubleInt),
		mapTestCase("doubleIntEmpty", []int{}, []int{}, doubleInt),
		mapTestCase("doubleIntNil", nil, nil, doubleInt),

		mapTestCase("halfFloat", list(1.0, 2, 3), list(0.5, 1, 1.5), halfFloat),
		mapTestCase("intToFloat", list(1, 2, 3), list(1.0, 2, 3), intToFloat),
		mapTestCase("intToHalfFloat", list(1, 2, 3), list(0.5, 1, 1.5), intToHalfFloat),

		PanicTestCases(func(s iter.Seq[int]) iter.Seq[int] {
			return Map(s, func(x int) int { return x })
		}),
	}.Run(t)
}

func pairUp[T, U any](a []T, b []U) []pair.Pair[T, U] {
	if len(a) != len(b) {
		panic("pairUp: slices must be the same length")
	}
	result := make([]pair.Pair[T, U], len(a))
	for i := range result {
		result[i] = pair.Of(a[i], b[i])
	}
	return result
}

func mapInTestCase[T, U, V any](name string, source []pair.Pair[T, U], want []V, f func(T, U) V) *SimpleTestCase[[]V] {
	source2 := MapOut(slices.Values(source), (pair.Pair[T, U]).Unpack)
	return SliceCollectTest(name, MapIn(source2, f), want)
}

func TestMapIn(t *testing.T) {
	add := func(a, b int) int { return a + b }
	sub := func(a, b int) int { return a - b }

	TestSuite{
		mapInTestCase("add", pairUp(list(1, 2, 3), list(4, 5, 6)), list(5, 7, 9), add),
		mapInTestCase("sub", pairUp(list(1, 2, 3), list(4, 5, 6)), list(-3, -3, -3), sub),

		mapInTestCase("addNil", nil, nil, add),
		mapInTestCase("subNil", nil, nil, sub),
		mapInTestCase("addEmpty", []pair.Pair[int, int]{}, []int{}, add),
		mapInTestCase("subEmpty", []pair.Pair[int, int]{}, []int{}, sub),

		PanicTestCases(func(s iter.Seq[pair.Pair[int, int]]) iter.Seq[int] {
			unpaired := MapOut(s, (pair.Pair[int, int]).Unpack)
			return MapIn(unpaired, add)
		}),
	}.Run(t)
}

func mapOutTestCase[T, V, W any](name string, source []T, want []pair.Pair[V, W], f func(T) (V, W)) *SimpleTestCase[[]pair.Pair[V, W]] {
	got := MapOut(slices.Values(source), f)
	return SliceCollectTest(name, MapIn(got, pair.Of), want)
}

func TestMapOut(t *testing.T) {
	isOdd := func(i int) (int, bool) { return i, i%2 == 1 }
	round := func(f float64) float64 { return math.Round(f*1e6) / 1e6 }
	sinCos := func(f float64) (float64, float64) { return round(math.Sin(f)), round(math.Cos(f)) }

	TestSuite{
		mapOutTestCase("isOdd", list(1, 2, 3), pairUp(list(1, 2, 3), list(true, false, true)), isOdd),
		mapOutTestCase("sinCos", list(0.0, math.Pi/2, math.Pi), pairUp(list(0.0, 1, 0), list(1.0, 0, -1)), sinCos),

		mapOutTestCase("isOddNil", nil, nil, isOdd),
		mapOutTestCase("isOddEmpty", []int{}, []pair.Pair[int, bool]{}, isOdd),
		mapOutTestCase("sinCosNil", nil, nil, sinCos),
		mapOutTestCase("sinCosEmpty", []float64{}, []pair.Pair[float64, float64]{}, sinCos),

		PanicTestCases(func(s iter.Seq[float64]) iter.Seq[pair.Pair[float64, float64]] {
			return MapIn(MapOut(s, sinCos), pair.Of)
		}),
	}.Run(t)
}
