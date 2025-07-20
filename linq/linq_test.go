package linq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -------------------------------------------------------------------
// 1 · Where → Select → ToSlice  (two transformers  ➜  Pipe3)
// -------------------------------------------------------------------
func TestStream_SelectWhere_ToSlice(t *testing.T) {
	ctx := t.Context()
	data := []int{1, 2, 3, 4, 5}

	got, err := Pipe3(
		FromSlice(ctx, data),                        // source
		Where(func(n int) bool { return n%2 == 1 }), // keep odds
		Select(func(n int) int { return n * n }),    // square each
		ToSlice[int](),                              // sink
	)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 9, 25}, got,
		"where→select chain should filter and transform correctly")
}

// -------------------------------------------------------------------
// 2 · Any / All / Count / First
// -------------------------------------------------------------------
func TestStream_Any_All_Count_First(t *testing.T) {
	ctx := t.Context()
	q := []int{1, 2, 3, 4}

	// ---- Any -------------------------------------------------------
	anyRes, err := Pipe1(
		FromSlice(ctx, q),
		Any(func(n int) bool { return n > 3 }), // any number > 3
	)
	assert.NoError(t, err)
	assert.True(t, anyRes)

	anyRes, err = Pipe1(
		FromSlice(ctx, q),
		Any(func(n int) bool { return n > 10 }), // any number > 10
	)
	assert.NoError(t, err)
	assert.False(t, anyRes)

	// ---- All -------------------------------------------------------
	allRes, err := Pipe1(
		FromSlice(ctx, q),
		All(func(n int) bool { return n < 5 }), // all numbers < 5
	)
	assert.NoError(t, err)
	assert.True(t, allRes)

	allRes, err = Pipe1(
		FromSlice(ctx, q),
		All(func(n int) bool { return n%2 == 0 }), // all numbers even
	)
	assert.NoError(t, err)
	assert.False(t, allRes)

	// ---- Count even numbers:  Where ➜ Count  (one transformer ➜ Pipe2)
	cnt, err := Pipe2(
		FromSlice(ctx, q),
		Where(func(n int) bool { return n%2 == 0 }),
		Count[int](),
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)

	// ---- First even number  (Where ➜ First)
	first, err := Pipe2(
		FromSlice(ctx, q),
		Where(func(n int) bool { return n%2 == 0 }),
		First[int](),
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, first)

	// ---- First > 10  (should not exist)
	_, err = Pipe2(
		FromSlice(ctx, q),
		Where(func(n int) bool { return n > 10 }),
		First[int](),
	)
	assert.Error(t, err)
	assert.EqualError(t, err, "stream is empty, no first element found")

	// ---- Total length (Count sink directly)
	total, err := Pipe1(
		FromSlice(ctx, q),
		Count[int](),
	)

	assert.NoError(t, err)
	assert.Equal(t, 4, total)
}

// -------------------------------------------------------------------
// 3 · Reverse strings (no transformer, call sink directly)
// -------------------------------------------------------------------
func TestStream_GenericWithStrings_Reverse(t *testing.T) {
	ctx := t.Context()

	got, err := Pipe1(
		FromSlice(ctx, []string{"a", "b", "c"}),
		Reverse[string](),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "b", "a"}, got)
}

func Test_Distinct(t *testing.T) {
	ctx := t.Context()

	type TestData struct {
		ID   int
		Name string
	}

	data := []TestData{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice"}, // duplicate by ID
	}

	got, err := Pipe2(
		FromSlice(ctx, data),
		Distinct(func(d TestData) int { return d.ID }),
		ToSlice[TestData](),
	)

	expected := []TestData{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, got, "Distinct should remove duplicates based on ID")
}

func Test_Flatten(t *testing.T) {
	ctx := t.Context()

	data := [][]int{
		{1, 2},
		{3, 4},
		{5},
	}

	got, err := Pipe2(
		FromSlice(ctx, data),
		Flatten[int](),
		ToSlice[int](),
	)

	expected := []int{1, 2, 3, 4, 5}

	assert.NoError(t, err)
	assert.Equal(t, expected, got, "Flatten should concatenate nested slices")
}
