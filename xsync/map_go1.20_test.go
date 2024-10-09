//go:build go1.20
// +build go1.20

package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapCASMethods(t *testing.T) {
	a := assert.New(t)

	var m Map[string, int]
	m.Store("one", 11)
	m.Store("two", 22)

	t.Log("Initial state")
	checkMap(t, &m, map[string]int{"one": 11, "two": 22})

	m.CompareAndSwap("one", 11, 1)
	t.Log("After CompareAndSwap of correct value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 22})

	m.CompareAndSwap("two", 3, 2)
	t.Log("After CompareAndSwap of incorrect value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 22})

	m.CompareAndDelete("two", 22)
	t.Log("After CompareAndDelete of correct value")
	checkMap(t, &m, map[string]int{"one": 1})

	m.CompareAndDelete("one", 11)
	t.Log("After CompareAndDelete of incorrect value")
	checkMap(t, &m, map[string]int{"one": 1})

	old, swap := m.Swap("two", 23)
	t.Log("After Swap of deleted/non-existant key")
	a.Zero(old, "Swap should return zero value for missing key")
	a.False(swap, "Swap should return false for missing key")
	checkMap(t, &m, map[string]int{"one": 1, "two": 23})

	old, swap = m.Swap("two", 2)
	t.Log("After Swap of existing key")
	a.Equal(23, old, "Swap should return previous value for existing key")
	a.True(swap, "Swap should return true for existing key")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2})
}
