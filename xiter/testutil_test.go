package xiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GenericTestCase interface {
	Name() string
	Run(t *testing.T)
}

type GenericTestCases []GenericTestCase

func (g GenericTestCases) Run(t *testing.T) {
	for _, tc := range g {
		if tcs, ok := tc.(GenericTestCases); ok {
			tcs.Run(t)
			continue
		}
		t.Run(tc.Name(), tc.Run)
	}
}

func (g GenericTestCases) Name() string {
	return "Generic test cases"
}

type SimpleTestCase[T any] struct {
	name       string
	generator  func(t *testing.T) T
	checks     []assert.ValueAssertionFunc
	args       []any
	panics     bool
	panicValue any
}

func SimpleTest[T any](name string, generator func(*testing.T) T) *SimpleTestCase[T] {
	return &SimpleTestCase[T]{
		name:      name,
		generator: generator,
	}
}

func (s *SimpleTestCase[T]) Panics() *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = nil
	return s
}

func (s *SimpleTestCase[T]) PanicsError(err error) *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = err
	return s
}

func (s *SimpleTestCase[T]) PanicsWith(panicValue any) *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = panicValue
	return s
}

func (s *SimpleTestCase[T]) Args(msgargs ...any) *SimpleTestCase[T] {
	s.args = msgargs
	return s
}

func (s *SimpleTestCase[T]) doChecks(t *testing.T, value T) {
	for _, check := range s.checks {
		check(t, value)
	}
}

func (s *SimpleTestCase[T]) Name() string {
	return s.name
}

func (s *SimpleTestCase[T]) Run(t *testing.T) {
	var value T
	if s.panics {
		if s.panicValue == nil {
			assert.Panics(t, func() {
				value = s.generator(t)
			}, s.args...)
		} else {
			if errValue, ok := s.panicValue.(error); ok {
				assert.PanicsWithError(t, errValue.Error(), func() {
					value = s.generator(t)
				}, s.args...)
			} else {
				assert.PanicsWithValue(t, s.panicValue, func() {
					value = s.generator(t)
				}, s.args...)
			}
		}
	} else {
		assert.NotPanics(t, func() {
			value = s.generator(t)
		}, s.args...)
	}
	for _, check := range s.checks {
		check(t, value)
	}
}

func (s *SimpleTestCase[T]) Value(vaf assert.ValueAssertionFunc) *SimpleTestCase[T] {
	s.checks = append(s.checks, vaf)
	return s
}

func (s *SimpleTestCase[T]) Compare(to T, caf assert.ComparisonAssertionFunc) *SimpleTestCase[T] {
	return s.Value(func(tt assert.TestingT, got interface{}, i2 ...interface{}) bool {
		return caf(tt, got, to, i2...)
	})
}

func CountUses(iter func(func(int) bool)) (func(func(int) bool), *int) {
	n := 0
	return func(yield func(int) bool) {
		n++
		iter(yield)
	}, &n
}

func IteratorCollectTest[T any](name string, it func(func(T) bool), want []T) *SimpleTestCase[[]T] {
	tc := SimpleTest(name, func(t *testing.T) []T {
		return sliceCollect(it)
	})
	if len(want) == 0 {
		return tc.Value(assert.Empty).Args("match empty slice")
	}
	return tc.Compare(want, assert.EqualValues).Args("match slice")
}

func PanicTestCases[T any](f func(func(func(T) bool)) func(func(T) bool)) GenericTestCases {
	return GenericTestCases{
		SimpleTest("nilIterPanic", func(t *testing.T) bool {
			var it func(func(T) bool)
			out := f(it)
			sliceCollect(out)
			return true
		}).Panics().Args("panics as expected"),
		SimpleTest("iterPanic", func(t *testing.T) bool {
			it := func(func(T) bool) { panic("no values") }
			out := f(it)
			sliceCollect(out)
			return true
		}).PanicsWith("no values").Args("panics as expected"),
	}
}
