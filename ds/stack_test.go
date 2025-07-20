package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper that drains the stack into a slice (top → bottom).
func drain[T any](st *Stack[T]) (out []T) {
	for !st.Empty() {
		v, _ := st.Pop()
		out = append(out, v)
	}
	return
}

func TestPushPopInt(t *testing.T) {
	st := NewStack[int]()

	// Push 1..3
	for i := 1; i <= 3; i++ {
		st.Push(i)
	}

	assert.Equal(t, 3, st.Len(), "length after pushes")

	// Pop should be 3,2,1
	got := drain(st)
	assert.Equal(t, []int{3, 2, 1}, got, "LIFO order")

	// Stack should now be empty
	_, ok := st.Pop()
	assert.False(t, ok, "pop on empty returns ok==false")
	assert.True(t, st.Empty(), "stack reports empty")
}

func TestPeek(t *testing.T) {
	st := NewStack[string]()
	st.Push("alpha")
	st.Push("beta")

	top, ok := st.Peek()
	assert.True(t, ok, "peek should succeed")
	assert.Equal(t, "beta", top, "peek shows top element")
	assert.Equal(t, 2, st.Len(), "peek leaves stack unchanged")
}

func TestGenericTypes(t *testing.T) {
	type user struct{ ID int }

	st := NewStack[user]()

	u := user{ID: 42}
	st.Push(u)

	got, ok := st.Pop()
	assert.True(t, ok)
	assert.Equal(t, u, got)
}
