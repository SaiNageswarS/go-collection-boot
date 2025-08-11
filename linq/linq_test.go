package linq

import (
	"context"
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

func Test_ForEach(t *testing.T) {
	ctx := t.Context()

	var sum int
	_, err := Pipe1(
		FromSlice(ctx, []int{1, 2, 3}),
		ForEach(func(n int) { sum += n }),
	)

	assert.NoError(t, err)
	assert.Equal(t, 6, sum, "ForEach should apply function to each element")
}

func Test_PipeUpstreamCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // upstream cancels the context immediately

	data := []int{1, 2, 3, 4, 5}
	got, err := Pipe3(
		FromSlice(ctx, data),
		Where(func(n int) bool { return n%2 == 1 }),
		Select(func(n int) int { return n * n }),
		ToSlice[int](),
	)

	assert.Error(t, err)
	assert.Equal(t, 0, len(got), "should return nil slice on immediate cancel")
}

func Test_trySend_ImmediateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel up-front

	ch := make(chan int)       // no receiver; send would block
	ok := trySend(ctx, ch, 42) // should notice ctx.Done first

	assert.False(t, ok, "trySend must return false when ctx is already cancelled")
}

func Test_SelectPar(t *testing.T) {
	ctx := t.Context()

	data := []int{1, 2, 3, 4, 5}

	// SelectPar should apply the function in parallel
	got, err := Pipe2(
		FromSlice(ctx, data),
		SelectPar(func(n int) int { return n * n }), // square each in parallel
		ToSlice[int](),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 4, 9, 16, 25}, got,
		"SelectPar should apply function to each element in parallel")
}

func Test_GroupBy(t *testing.T) {
	ctx := t.Context()

	type sample struct {
		employeeID string
		department string
		name       string
	}
	data := []sample{
		{employeeID: "1", department: "HR", name: "Alice"},
		{employeeID: "2", department: "IT", name: "Bob"},
		{employeeID: "3", department: "HR", name: "Charlie"},
		{employeeID: "4", department: "IT", name: "David"},
		{employeeID: "5", department: "Finance", name: "Eve"},
		{employeeID: "6", department: "IT", name: "Frank"},
	}

	// GroupBy should group numbers by even/odd
	got, err := Pipe3(
		FromSlice(ctx, data),
		GroupBy(func(s sample) string { return s.department }), // group by department
		Select(func(g []sample) int { return len(g) }),         // count employees in each group
		ToSlice[int](), // collect results
	)

	assert.NoError(t, err)
	expected := []int{2, 3, 1} // HR: 2, IT: 3, Finance: 1
	// The order of groups is not guaranteed, so we can only check the counts
	// Check if we have the expected counts
	assert.Len(t, got, 3, "should have 3 groups")
	assert.Contains(t, got, expected[0], "should have 2 employees in HR")
	assert.Contains(t, got, expected[1], "should have 2 employees in IT")
	assert.Contains(t, got, expected[2], "should have 1 employee in Finance")
}

func Test_GroupBy_Any(t *testing.T) {
	ctx := t.Context()

	type sample struct {
		employeeID string
		department string
		name       string
	}
	data := []sample{
		{employeeID: "1", department: "HR", name: "Alice"},
		{employeeID: "2", department: "IT", name: "Bob"},
		{employeeID: "3", department: "HR", name: "Charlie"},
		{employeeID: "4", department: "IT", name: "David"},
		{employeeID: "5", department: "HR", name: "Eve"},
		{employeeID: "6", department: "IT", name: "Frank"},
		{employeeID: "16", department: "Finance", name: "Paul"},
		{employeeID: "17", department: "Finance", name: "Quinn"},
		{employeeID: "18", department: "Finance", name: "Rita"},
		{employeeID: "19", department: "Sales", name: "Sam"},
	}

	// GroupBy should group by department and check if any group has more than 2 employees
	// This will end the pipeline early if any group has more than 2 employees
	got, err := Pipe2(
		FromSlice(ctx, data),
		GroupBy(func(s sample) string { return s.department }),
		Any(func(g []sample) bool {
			return len(g) > 2 // check if any group has two employees
		}),
	)

	assert.NoError(t, err)
	assert.True(t, got, "should return true since IT, HR and Finance departments have more than 2 employees")
}

func Test_EarlyCancel_PropagatesThroughTransformers(t *testing.T) {
	ctx := t.Context()

	input := [][]int{{1, 2}, {3}} // at least two unique ints so Distinct tries a 2nd send

	// Pipeline:  FromSlice ➜ Where ➜ Select ➜ Distinct ➜ First
	// `First` cancels after emitting 1; transformers then hit ctx.Done().
	first, err := Pipe5(
		fromUnbufferedSlice(ctx, input),
		Flatten[int](),
		Where(func(n int) bool { return n > 0 }), // pass all
		Select(func(n int) int { return n }),     // identity
		Distinct(func(n int) int { return n }),   // dedupe
		First[int](),                             // sink – cancels early
	)

	assert.NoError(t, err)
	assert.Equal(t, 1, first)
}

func fromUnbufferedSlice[T any](parent context.Context, src []T) Stream[T] {
	ctx, cancel := context.WithCancel(parent)
	out := make(chan T) // unbuffered

	go func() {
		defer close(out)
		for _, v := range src {
			if !trySend(ctx, out, v) { // stop if downstream cancelled
				return
			}
		}
	}()

	return NewStream(ctx, out, cancel, 0)
}
