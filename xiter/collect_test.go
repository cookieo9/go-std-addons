package xiter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func collectTestCase[T, Accum any](name string, source []T, start Accum, f func(Accum, T) Accum, want Accum) GenericTestCase {
	testCase := SimpleTest(name, func(t *testing.T) Accum {
		return Collect(sliceValues(source), start, f)
	})
	if rv := reflect.ValueOf(want); rv.Kind() == reflect.Slice && rv.Len() == 0 {
		return testCase.Value(assert.Empty).Args("empty result")
	}
	return testCase.Compare(want, assert.EqualValues).Args("same result")
}

func TestCollect(t *testing.T) {
	sumAccum := func(accum int, t int) int { return accum + t }
	collectAccum := func(accum []int, t int) []int { return append(accum, t) }
	GenericTestCases{
		collectTestCase("emptySum", []int{}, 0, sumAccum, 0),
		collectTestCase("emptySumNil", nil, 0, sumAccum, 0),
		collectTestCase("emptySumStart10", []int{}, 10, sumAccum, 10),

		collectTestCase("emptyCollect", []int{}, []int{}, collectAccum, []int{}),
		collectTestCase("emptyCollectNil", nil, []int{}, collectAccum, []int{}),
		collectTestCase("emptyCollectNilNil", nil, nil, collectAccum, []int{}),
		collectTestCase("emptyCollectStart10", []int{}, []int{10}, collectAccum, []int{10}),

		collectTestCase("someSum", []int{1, 2, 3}, 0, sumAccum, 6),
		collectTestCase("someSumStart10", []int{1, 2, 3}, 10, sumAccum, 16),

		collectTestCase("someCollect", []int{1, 2, 3}, []int{}, collectAccum, []int{1, 2, 3}),
		collectTestCase("someCollectStart10", []int{1, 2, 3}, []int{10}, collectAccum, []int{10, 1, 2, 3}),

		PanicTestCases(func(f func(func(int) bool)) func(func(int) bool) {
			Collect(f, 0, sumAccum)
			return f
		}),
	}.Run(t)
}
