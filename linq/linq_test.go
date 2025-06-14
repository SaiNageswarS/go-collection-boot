package linq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_Chaining_SelectWhere_ToSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	got := From(data).
		Where(func(n int) bool { return n%2 == 1 }). // keep odds
		Select(func(n int) int { return n * n }).    // square each
		ToSlice()

	want := []int{1, 9, 25}
	assert.Equal(t, want, got, "where→select chain should filter and transform correctly")
}

func TestQuery_Any_All_Count_First(t *testing.T) {
	q := From([]int{1, 2, 3, 4})

	assert.True(t, q.Any(func(n int) bool { return n > 3 }))
	assert.False(t, q.Any(func(n int) bool { return n > 10 }))

	assert.True(t, q.All(func(n int) bool { return n < 5 }))
	assert.False(t, q.All(func(n int) bool { return n%2 == 0 }))

	assert.Equal(t, 2, q.Count(func(n int) bool { return n%2 == 0 }))

	first, ok := q.First(func(n int) bool { return n%2 == 0 })
	assert.True(t, ok)
	assert.Equal(t, 2, first)

	_, ok = q.First(func(n int) bool { return n > 10 })
	assert.False(t, ok)
}

func TestQuery_Distinct_Reverse_Len(t *testing.T) {
	q := From([]int{1, 2, 2, 3, 3, 3})

	distinct := q.Distinct().ToSlice()
	assert.Equal(t, []int{1, 2, 3}, distinct)

	reversed := q.Reverse().ToSlice()
	assert.Equal(t, []int{3, 3, 3, 2, 2, 1}, reversed)

	assert.Equal(t, 6, q.Len())
}

func TestQuery_SetOps(t *testing.T) {
	a := From([]int{1, 2, 3})
	b := From([]int{3, 4})

	union := a.Union(b).ToSlice()
	assert.Equal(t, []int{1, 2, 3, 4}, union, "union keeps order of first query then uniques from second")

	inter := a.Intersection(b).ToSlice()
	assert.Equal(t, []int{3}, inter)

	diff := a.Difference(b).ToSlice()
	assert.Equal(t, []int{1, 2}, diff)
}

func TestQuery_GenericWithStrings(t *testing.T) {
	q := From([]string{"a", "b", "a", "c"}).
		Distinct().
		Reverse().
		ToSlice()

	assert.Equal(t, []string{"c", "b", "a"}, q, "distinct should remove duplicates")
}
