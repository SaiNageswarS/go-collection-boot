package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndContains(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)

	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.False(t, s.Contains(4))
}

func TestAddAll(t *testing.T) {
	s := NewSet[string]()
	s.AddAll([]string{"a", "b", "c"})

	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.True(t, s.Contains("c"))
}

func TestRemove(t *testing.T) {
	s := NewSet[int]()
	s.Add(5)
	s.Remove(5)

	assert.False(t, s.Contains(5))
}

func TestContainsAll(t *testing.T) {
	s := NewSet[int]()
	s.AddAll([]int{10, 20, 30})

	assert.True(t, s.ContainsAll(10, 20))
	assert.False(t, s.ContainsAll(10, 40))
}

func TestToSliceAndLen(t *testing.T) {
	s := NewSet[string]()
	s.Add("x", "y")

	slice := s.ToSlice()
	assert.Len(t, slice, 2)
	assert.Equal(t, 2, s.Len())
}

func TestClear(t *testing.T) {
	s := NewSet[int]()
	s.AddAll([]int{1, 2, 3})
	s.Clear()

	assert.Equal(t, 0, s.Len())
	assert.True(t, s.IsEmpty())
}

func TestClone(t *testing.T) {
	s1 := NewSet[int]()
	s1.Add(1, 2)
	s2 := s1.Clone()
	s2.Add(3)

	assert.True(t, s1.Contains(1))
	assert.False(t, s1.Contains(3))
	assert.True(t, s2.Contains(3))
}

func TestUnion(t *testing.T) {
	a := FromSlice([]string{"a", "b"})
	b := FromSlice([]string{"b", "c"})

	u := a.Union(b)

	assert.True(t, u.Contains("a"))
	assert.True(t, u.Contains("b"))
	assert.True(t, u.Contains("c"))
	assert.Equal(t, 3, u.Len())
}

func TestIntersection(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := FromSlice([]int{2, 3, 4})

	inter := a.Intersection(b)

	assert.True(t, inter.Contains(2))
	assert.True(t, inter.Contains(3))
	assert.False(t, inter.Contains(1))
	assert.False(t, inter.Contains(4))
	assert.Equal(t, 2, inter.Len())
}

func TestDifference(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := FromSlice([]int{3, 4, 5})

	diff := a.Difference(b)

	assert.True(t, diff.Contains(1))
	assert.True(t, diff.Contains(2))
	assert.False(t, diff.Contains(3))
}
