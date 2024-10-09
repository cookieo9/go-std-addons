package xsync

import (
	"testing"

	"github.com/cookieo9/go-std-addons/pair"
	"github.com/stretchr/testify/assert"
)

func TestMapInvariants(t *testing.T) {
	a := assert.New(t)
	var m Map[string, int]

	t.Log("zero state")
	a.Equal(0, MapLen(&m), "Length should be zero initially")
	checkMap(t, &m, map[string]int{})

	x, loaded := m.LoadOrStore("one", 1)
	t.Log("After first LoadOrStore")
	a.False(loaded, "First LoadOrStore should return false for loaded")
	a.Equal(x, 1, "First LoadOrStore should return assigned value")
	checkMap(t, &m, map[string]int{"one": 1})

	x, loaded = m.LoadOrStore("one", 100)
	t.Log("After LoadOrStore with existing key")
	a.True(loaded, "Second LoadOrStore should return true for loaded")
	a.Equal(x, 1, "Second LoadOrStore should return old value")
	checkMap(t, &m, map[string]int{"one": 1})

	x, loaded = m.LoadAndDelete("two")
	t.Log("After LoadAndDelete with missing key")
	a.False(loaded, "LoadAndDelete should return false for missing key")
	a.Zero(x, "LoadAndDelete should return zero value for missing key")
	checkMap(t, &m, map[string]int{"one": 1})

	m.Store("two", 22)
	t.Log("After store of new key 'two")
	checkMap(t, &m, map[string]int{"one": 1, "two": 22})

	x, loaded = m.Load("two")
	t.Log("After Load of existing key")
	a.True(loaded, "Load should return true for existing key")
	a.Equal(x, 22, "Load should return stored value")
	checkMap(t, &m, map[string]int{"one": 1, "two": 22})

	m.Delete("two")
	t.Log("After Delete of existing key")
	checkMap(t, &m, map[string]int{"one": 1})

	x, loaded = m.Load("two")
	t.Log("After Load of deleted key")
	a.False(loaded, "Load should return false for deleted key")
	a.Zero(x, "Load should return zero value for deleted key")
	checkMap(t, &m, map[string]int{"one": 1})

	m.Store("three", 3)
	t.Log("After Store of new key 'three'")
	checkMap(t, &m, map[string]int{"one": 1, "three": 3})

	x, loaded = m.LoadAndDelete("three")
	t.Log("After LoadAndDelete of existing key 'three'")
	a.True(loaded, "LoadAndDelete should return true for existing key")
	a.Equal(x, 3, "LoadAndDelete should return stored value")
	checkMap(t, &m, map[string]int{"one": 1})

	x, loaded = m.Load("three")
	t.Log("After Load of deleted key 'three'")
	a.False(loaded, "Load should return false for missing key")
	a.Zero(x, "Load should return zero value for missing key")
	checkMap(t, &m, map[string]int{"one": 1})

	m.Store("four", 4)
	m.Store("five", 5)
	t.Log("After multiple Stores (with new keys)")
	checkMap(t, &m, map[string]int{"one": 1, "four": 4, "five": 5})
	a.Equal(pair.Of(1, true), pair.Of(m.Load("one")), "Load should return 1 for value for 'one'")
	a.Equal(pair.Of(4, true), pair.Of(m.Load("four")), "Load should return 4 for value for 'four'")
	a.Equal(pair.Of(5, true), pair.Of(m.Load("five")), "Load should return 5 for value for 'five'")
}
