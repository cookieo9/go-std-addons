package pipe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cookieo9/go-std-addons/xiter/pipe"
)

func TestPipeline(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	want := []float64{1.1, 3.3}
	isOdd := func(x int) bool { return x%2 == 1 }

	got, err := pipe.ProcessSlice[float64](data,
		pipe.Filter(isOdd),
		pipe.Map(func(x int) float64 { return float64(x) * 1.1 }),
		pipe.While(func(f float64) bool { return f < 4 }),
		pipe.Materialize[float64](),
	)
	require.NoError(t, err)
	assert.InEpsilonSlice(t, want, got, 0.0000001, "same sequence")
}

func TestJoins(t *testing.T) {
	double := pipe.Map[int](func(x int) int { return x * 2 })
	odd := pipe.Filter(func(x int) bool { return x%2 == 1 })
	halfFloat := pipe.Map(func(x float64) float64 { return float64(x) / 2 })
	testCases := []struct {
		name         string
		stages       []pipe.Processor
		joinErr      assert.ErrorAssertionFunc
		panicProcess bool
	}{
		{name: "no stages"},
		{name: "double", stages: []pipe.Processor{double}},
		{name: "odd, double", stages: []pipe.Processor{odd, double}},
		{name: "double, halfFloat", stages: []pipe.Processor{double, halfFloat}, joinErr: assert.Error},
		{name: "halfFloat, double", stages: []pipe.Processor{halfFloat, double}, joinErr: assert.Error},
		{name: "halfFloat", stages: []pipe.Processor{halfFloat}, panicProcess: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := pipe.TryJoin(tc.stages...)
			if tc.joinErr != nil {
				tc.joinErr(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			if tc.panicProcess {
				require.Panics(t, func() {
					pipe.ProcessSlice[int]([]int{1, 2, 3}, p)
				})
				return
			}
			require.NotPanics(t, func() {
				pipe.ProcessSlice[int]([]int{1, 2, 3}, p)
			})
		})
	}
}
