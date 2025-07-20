package ds

type Stack[T any] struct {
	data []T
}

// NewStack returns an empty stack with optional initial capacity.
func NewStack[T any](capHint ...int) *Stack[T] {
	capacity := 0
	if len(capHint) > 0 && capHint[0] > 0 {
		capacity = capHint[0]
	}
	return &Stack[T]{data: make([]T, 0, capacity)}
}

// Push adds v to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

// Pop removes the top element and returns it.
// The second return value is false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	idx := len(s.data) - 1
	v := s.data[idx]
	// Avoid memory leaks for reference types.
	s.data[idx] = zero
	s.data = s.data[:idx]
	return v, true
}

// Peek returns—but does not remove—the top element.
// The second return value is false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Len reports the current number of elements.
func (s *Stack[T]) Len() int { return len(s.data) }

// Empty is a convenience wrapper around Len().
func (s *Stack[T]) Empty() bool { return len(s.data) == 0 }

// Clear discards all elements (retains allocated capacity).
func (s *Stack[T]) Clear() {
	var zero T
	for i := range s.data {
		s.data[i] = zero // help GC for reference types
	}
	s.data = s.data[:0]
}
