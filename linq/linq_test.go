package linq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	result := Where(data, func(n int) bool { return n%2 == 0 })

	assert.Equal(t, []int{2, 4}, result)
}

func TestSelect(t *testing.T) {
	data := []string{"a", "bb", "ccc"}
	result := Select(data, func(s string) int { return len(s) })

	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestAny(t *testing.T) {
	data := []int{1, 3, 5}
	assert.False(t, Any(data, func(n int) bool { return n%2 == 0 }))
	assert.True(t, Any(data, func(n int) bool { return n == 3 }))
}

func TestAll(t *testing.T) {
	data := []int{2, 4, 6}
	assert.True(t, All(data, func(n int) bool { return n%2 == 0 }))
	assert.False(t, All(data, func(n int) bool { return n > 2 }))
}

func TestFirst(t *testing.T) {
	data := []string{"dog", "cat", "cow"}
	item, ok := First(data, func(s string) bool { return strings.HasPrefix(s, "c") })

	assert.True(t, ok)
	assert.Equal(t, "cat", item)

	_, ok = First(data, func(s string) bool { return s == "elephant" })
	assert.False(t, ok)
}

func TestCount(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	count := Count(data, func(n int) bool { return n%2 == 0 })

	assert.Equal(t, 3, count)
}

func TestDistinct(t *testing.T) {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	result := Distinct(data)

	assert.ElementsMatch(t, []string{"apple", "banana", "cherry"}, result)
}

func TestReverse(t *testing.T) {
	data := []int{1, 2, 3}
	result := Reverse(data)

	assert.Equal(t, []int{3, 2, 1}, result)

	// original should remain unchanged
	assert.Equal(t, []int{1, 2, 3}, data)
}
