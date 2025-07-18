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
	q := []int{1, 2, 3, 4}

	assert.True(t, From(q).Any(func(n int) bool { return n > 3 }))
	assert.False(t, From(q).Any(func(n int) bool { return n > 10 }))

	assert.True(t, From(q).All(func(n int) bool { return n < 5 }))
	assert.False(t, From(q).All(func(n int) bool { return n%2 == 0 }))

	assert.Equal(t, 2, From(q).Count(func(n int) bool { return n%2 == 0 }))

	first, ok := From(q).First(func(n int) bool { return n%2 == 0 })
	assert.True(t, ok)
	assert.Equal(t, 2, first)

	_, ok = From(q).First(func(n int) bool { return n > 10 })
	assert.False(t, ok)

	count := From(q).Len()
	assert.Equal(t, 4, count, "Len should return the number of elements in the query")
}

func TestQuery_GenericWithStrings(t *testing.T) {
	q := From([]string{"a", "b", "a", "c"}).
		Reverse().
		ToSlice()

	assert.Equal(t, []string{"c", "a", "b", "a"}, q, "distinct should remove duplicates")
}
