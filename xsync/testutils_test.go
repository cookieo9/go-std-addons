package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkMap(t *testing.T, m *Map[string, int], want map[string]int) {
	t.Helper()

	// Check length
	assert.Equal(t, len(want), MapLen(m), "Length should be %d", len(want))

	got := make(map[string]int, len(want))
	m.Range(func(k string, v int) bool {
		got[k] = v
		return true
	})

	assert.Exactly(t, want, got, "Map should contain expected items")
}
