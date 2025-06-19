package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinHeap(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	h := NewMinHeap(less)

	h.Push(5)
	h.Push(3)
	h.Push(8)
	h.Push(1)

	assert.Equal(t, 4, h.Len(), "Heap length should be 4 after 4 pushes")
	min, ok := h.Peek()
	assert.True(t, ok, "Peek should return true")
	assert.Equal(t, 1, min, "Peek should return the minimum element")

	min, ok = h.Pop()
	assert.True(t, ok, "Pop should return true")
	assert.Equal(t, 1, min, "Pop should return the minimum element")
	assert.Equal(t, 3, h.Len(), "Heap length should be 3 after one pop")

	h.Push(2)
	assert.Equal(t, 4, h.Len(), "Heap length should be 4 after pushing 2")

	min, ok = h.Pop()
	assert.True(t, ok, "Pop should return true")
	assert.Equal(t, 2, min, "Pop should return the next minimum element")

	min, ok = h.Pop()
	assert.True(t, ok, "Pop should return true")
	assert.Equal(t, 3, min, "Pop should return the next minimum element")

	expectedRemaining := []int{5, 8}
	for idx, v := range h.ToSlice() {
		assert.Equal(t, expectedRemaining[idx], v, "Remaining elements should be in correct order")
	}
}
