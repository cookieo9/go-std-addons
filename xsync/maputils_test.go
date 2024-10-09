package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapClear(t *testing.T) {
	// Load initial values
	var m Map[string, int]
	m.Store("one", 1)
	m.Store("two", 2)
	m.Store("three", 3)

	t.Log("After initial values are loaded")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2, "three": 3})

	MapClear(&m)
	t.Log("After Clear")
	checkMap(t, &m, map[string]int{})
}

func TestMapCache(t *testing.T) {
	var m Map[string, int]
	m.Store("one", 1)
	m.Store("two", 2)

	t.Log("After initial values are loaded")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2})

	x, loaded := Cache(&m, "one", func() int { return 10 })
	t.Log("After Cache of existing key")
	assert.True(t, loaded, "Cache should return true for existing key")
	assert.Equal(t, x, 1, "Cache should return stored value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2})

	x, loaded = Cache(&m, "one-hundred", func() int { return 100 })
	t.Log("After Cache of missing key")
	assert.False(t, loaded, "Cache should return false for missing key")
	assert.Equal(t, x, 100, "Cache should return computed value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2, "one-hundred": 100})

	x, loaded = Cache(&m, "one-hundred", func() int { return 200 })
	t.Log("After Cache of item previously computed")
	assert.True(t, loaded, "Cache should return true for previously computed key")
	assert.Equal(t, x, 100, "Cache should return previously computed value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 2, "one-hundred": 100})
}
