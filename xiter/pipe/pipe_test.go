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

	got, err := pipe.ProcessSlice[float64](data, pipe.Pipeline(
		pipe.Filter(isOdd),
		pipe.Map(func(x int) float64 { return float64(x) * 1.1 }),
		pipe.While(func(f float64) bool { return f < 4 }),
		pipe.Materialize[float64](),
	))
	require.NoError(t, err)
	assert.InEpsilonSlice(t, want, got, 0.0000001, "same sequence")
}
