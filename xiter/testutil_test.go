package xiter

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase is an interface that defines the contract for a test case.
type TestCase interface {
	Name() string
	Run(t *testing.T)
}

// TestSuite is a collection of test cases that acts as a single test case.
type TestSuite []TestCase

// Run executes all the test cases in the TestSuite. If a test case is itself a
// TestSuite, it recursively runs the nested test cases.
func (g TestSuite) Run(t *testing.T) {
	for _, tc := range g {
		if tcs, ok := tc.(TestSuite); ok {
			tcs.Run(t)
			continue
		}
		t.Run(tc.Name(), tc.Run)
	}
}

// Name returns the name of the TestSuite to comply with the TestCase interface.
// It will never actually be called by a containing TestSuite.
func (g TestSuite) Name() string {
	return "Generic test cases"
}

// SimpleTestCase is a struct that represents a single modular test case. It's
// main method of operation is to call a generator function to produce a value
// to be checked (of type T), and checks are added via methods that accept
// assertions to call on it. The assertions are expected to be from the
// [github.com/stretchr/testify/assert] package.
type SimpleTestCase[T any] struct {
	name       string
	generator  func(t *testing.T) T
	checks     []assert.ValueAssertionFunc
	args       []any
	panics     bool
	panicValue any
}

// SimpleTest creates a new SimpleTestCase for the given generator function.
// The name parameter is used to identify the test case.
func SimpleTest[T any](name string, generator func(*testing.T) T) *SimpleTestCase[T] {
	return &SimpleTestCase[T]{
		name:      name,
		generator: generator,
	}
}

// Panics sets the SimpleTestCase to expect the generator function to panic.
// It will not check the value of the panic itself, only that it did. If the
// generator function does not panic, the test will fail.
func (s *SimpleTestCase[T]) Panics() *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = nil
	return s
}

// PanicsError sets the SimpleTestCase to expect the generator function to
// panic with the provided error value. If the generator function does not
// panic or the value of the panic does not match the provided error value,
// the test will fail. Matching errors is done via error string comparison.
func (s *SimpleTestCase[T]) PanicsError(err error) *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = err
	return s
}

// PanicsWith sets the SimpleTestCase to expect the generator function to
// panic with the provided panic value. If the generator function does not
// panic or the value of the panic does not match the provided panic value,
// the test will fail.
func (s *SimpleTestCase[T]) PanicsWith(panicValue any) *SimpleTestCase[T] {
	s.panics = true
	s.panicValue = panicValue
	return s
}

// Args sets the message arguments to be passed to assertions when they are
// checked. This is useful for adding context to the assertion messages. All
// assertions that accept message arguments will use the same arguments.
func (s *SimpleTestCase[T]) Args(msgargs ...any) *SimpleTestCase[T] {
	s.args = msgargs
	return s
}

// Name returns the name of the SimpleTestCase.
func (s *SimpleTestCase[T]) Name() string {
	return s.name
}

// Run executes the test case by running the generator function and checking the
// result. The test will fail if any of the checks fail, or if a panic / lack of
// panic doesn't match the expected behavior when running the generator.
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

// Value adds a value assertion function to the test case. The provided function
// will be called with the generated value, and the test will fail if the
// assertion fails.
func (s *SimpleTestCase[T]) Value(vaf assert.ValueAssertionFunc) *SimpleTestCase[T] {
	s.checks = append(s.checks, vaf)
	return s
}

// Compare adds a comparison assertion function to the test case. The provided
// function will be called with the generated value and the expected value, and
// the test will fail if the assertion fails.
func (s *SimpleTestCase[T]) Compare(to T, caf assert.ComparisonAssertionFunc) *SimpleTestCase[T] {
	return s.Value(func(tt assert.TestingT, got interface{}, i2 ...interface{}) bool {
		return caf(tt, got, to, i2...)
	})
}

// CountUses returns a new iterator that wraps the given iterator `it`, and also
// returns a pointer to an integer that counts the number of times the iterator
// is consumed.
func CountUses[T any](it iter.Seq[T]) (iter.Seq[T], *int) {
	n := 0
	return func(yield func(T) bool) {
		n++
		it(yield)
	}, &n
}

// SliceCollectTest is a test utility function that creates a SimpleTestCase
// meant to test the values generated by an iter.Seq[T] by collecting the values
// into a slice as the generator, and adding simple assertions based on an
// expected slice.
func SliceCollectTest[T any](name string, it iter.Seq[T], want []T) *SimpleTestCase[[]T] {
	tc := SimpleTest(name, func(t *testing.T) []T {
		return slices.Collect(it)
	})
	if len(want) == 0 {
		return tc.Value(assert.Empty).Args("match empty slice")
	}
	return tc.Compare(want, assert.EqualValues).Args("match slice")
}

// PanicTestCases creates a TestSuite that contains test cases for verifying
// that a given iterator processing function panics as expected when provided
// with invalid input.
//
// The TestSuite contains two test cases:
//
//  1. "nilIterPanic" - Verifies that the function panics when provided with a nil
//     iterator.
//  2. "iterPanic" - Verifies that the function panics when provided with an
//     iterator that panics.
//
// Most iterator processing functions can be tested for these situations.
func PanicTestCases[T, U any](f func(iter.Seq[T]) iter.Seq[U]) TestSuite {
	return TestSuite{
		SimpleTest("nilIterPanic", func(t *testing.T) bool {
			var it iter.Seq[T]
			out := f(it)
			slices.Collect(out)
			return true
		}).Panics().Args("panics as expected"),
		SimpleTest("iterPanic", func(t *testing.T) bool {
			it := func(func(T) bool) { panic("no values") }
			out := f(it)
			slices.Collect(out)
			return true
		}).PanicsWith("no values").Args("panics as expected"),
	}
}
