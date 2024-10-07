package xerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatch(t *testing.T) {
	testErr := errors.New("test")

	t.Run("panicWithError", func(t *testing.T) {
		f := func() {
			panic(testErr)
		}
		assert.ErrorIs(t, Catch(f), testErr, "error should be returned")
	})

	t.Run("panicWithNonError", func(t *testing.T) {
		assert.Panics(t, func() {
			f := func() {
				panic("not an error")
			}
			Catch(f)
		}, "panic should be re-thrown")
	})

	t.Run("noPanic", func(t *testing.T) {
		f := func() {
			// no panic
		}
		assert.NoError(t, Catch(f), "error should be nil")
	})
}

func TestCatchValue(t *testing.T) {
	testErr := errors.New("test")
	testValue := 42

	t.Run("panicWithError", func(t *testing.T) {
		f := func() int {
			panic(testErr)
		}
		value, err := CatchValue(f)
		assert.ErrorIs(t, err, testErr, "error should be returned")
		assert.Zero(t, value, "value should be zero")
	})

	t.Run("panicWithNonError", func(t *testing.T) {
		assert.Panics(t, func() {
			f := func() int {
				panic("not an error")
			}
			CatchValue(f)
		}, "panic should be re-thrown")
	})

	t.Run("noPanic", func(t *testing.T) {
		f := func() int {
			return testValue
		}
		value, err := CatchValue(f)
		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, testValue, value, "value should be returned")
	})
}
