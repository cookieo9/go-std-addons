package xiter

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountOverflowUnderflow(t *testing.T) {
	t.Run("OverflowUint8", func(t *testing.T) {
		iter := Limit(Count(uint8(250)), 10)
		s := slices.Collect(iter)
		assert.Equal(t, []uint8{250, 251, 252, 253, 254, 255, 0, 1, 2, 3}, s)
	})

	t.Run("OverflowInt8", func(t *testing.T) {
		iter := Limit(Count(int8(120)), 10)
		s := slices.Collect(iter)
		assert.Equal(t, []int8{120, 121, 122, 123, 124, 125, 126, 127, -128, -127}, s)
	})

	t.Run("UnderflowUint8", func(t *testing.T) {
		iter := Limit(CountDown(uint8(1), 1), 10)
		s := slices.Collect(iter)
		assert.Equal(t, []uint8{1, 0, 255, 254, 253, 252, 251, 250, 249, 248}, s)
	})

	t.Run("UnderflowInt8", func(t *testing.T) {
		iter := Limit(CountDown(int8(-125), 1), 10)
		s := slices.Collect(iter)
		assert.Equal(t, []int8{-125, -126, -127, -128, 127, 126, 125, 124, 123, 122}, s)
	})
}

func TestRanges(t *testing.T) {
	TestSuite{
		SliceCollectTest("Range(0,10)", Range(0, 10), []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
		SliceCollectTest("Range(1,10)", Range(1, 10), []int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		SliceCollectTest("Range(10,0)", Range(10, 0), []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),
		SliceCollectTest("Range(10,10)", Range(10, 10), []int{}),
		SliceCollectTest("Range(10,11)", Range(10, 11), []int{10}),
		SliceCollectTest("RangeBy(0,10,2)", RangeBy(0, 10, 2), []int{0, 2, 4, 6, 8}),
		SliceCollectTest("RangeBy(10,0,2)", RangeBy(10, 0, 2), []int{10, 8, 6, 4, 2}),
		SliceCollectTest("RangeBy(-5,5,3)", RangeBy(-5, 5, 3), []int{-5, -2, 1, 4}),

		SliceCollectTest("Range(0,10)float64", Range[float64](0, 10), []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
		SliceCollectTest("Range(1,10)float64", Range[float64](1, 10), []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		SliceCollectTest("Range(10,0)float64", Range[float64](10, 0), []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),
		SliceCollectTest("Range(10,10)float64", Range[float64](10, 10), []float64{}),
		SliceCollectTest("Range(10,11)float64", Range[float64](10, 11), []float64{10}),
		SliceCollectTest("RangeBy(0,10,2.5)float64", RangeBy[float64](0, 10, 2.5), []float64{0, 2.5, 5, 7.5}),
		SliceCollectTest("RangeBy(10,0,2.5)float64", RangeBy[float64](10, 0, 2.5), []float64{10, 7.5, 5, 2.5}),
		SliceCollectTest("RangeBy(-5,5,3.5)float64", RangeBy[float64](-5, 5, 3.5), []float64{-5, -1.5, 2}),
	}.Run(t)
}

func TestCountPanics(t *testing.T) {
	t.Run("Count", func(t *testing.T) {
		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Limit(Count(0), 10)
		}).Run(t)
	})
	t.Run("CountUp", func(t *testing.T) {
		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Limit(CountUp(0, 1), 10)
		}).Run(t)
	})
	t.Run("CountDown", func(t *testing.T) {
		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Limit(CountDown(0, 1), 10)
		}).Run(t)
	})

	t.Run("Range", func(t *testing.T) {
		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return Range(0, 10)
		}).Run(t)
	})

	t.Run("RangeBy", func(t *testing.T) {
		PanicTestCases(func(it iter.Seq[int]) iter.Seq[int] {
			slices.Collect(it)
			return RangeBy(0, 10, 1)
		}).Run(t)
	})
}
