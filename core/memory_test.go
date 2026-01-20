package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory(t *testing.T) {
	mem := NewMemory[string](3)

	mem.Push("first")
	assert.True(t, mem.Has("first"), "Has first")

	mem.Push("second")
	assert.True(t, mem.Has("first"), "Has first")
	assert.True(t, mem.Has("second"), "Has second")

	mem.Push("third")
	assert.True(t, mem.Has("first"), "Has first")
	assert.True(t, mem.Has("second"), "Has second")
	assert.True(t, mem.Has("third"), "Has third")

	mem.Push("fourth")
	assert.False(t, mem.Has("first"), "Doesn't have first")
	assert.True(t, mem.Has("second"), "Has second")
	assert.True(t, mem.Has("third"), "Has third")
	assert.True(t, mem.Has("fourth"), "Has fourth")
}
