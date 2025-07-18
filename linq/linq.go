package linq

// Query represents a lazy stream of T values.
type Query[T any] struct {
	C <-chan T
}

// From turns a slice into a Query that emits each element.
func From[T any](items []T) Query[T] {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range items {
			ch <- v
		}
	}()
	return Query[T]{C: ch}
}

// Where filters elements by predicate.
func (q Query[T]) Where(pred func(T) bool) Query[T] {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range q.C {
			if pred(v) {
				out <- v
			}
		}
	}()
	return Query[T]{C: out}
}

// Select maps each element through mapper.
func (q Query[T]) Select(mapper func(T) T) Query[T] {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range q.C {
			out <- mapper(v)
		}
	}()
	return Query[T]{C: out}
}

// ToSlice collects all remaining elements into a slice.
func (q Query[T]) ToSlice() []T {
	var result []T
	for v := range q.C {
		result = append(result, v)
	}
	return result
}

// Any returns true if predicate holds for any element (consumes the stream).
func (q Query[T]) Any(pred func(T) bool) bool {
	for v := range q.C {
		if pred(v) {
			return true
		}
	}
	return false
}

// All returns true if predicate holds for every element.
func (q Query[T]) All(pred func(T) bool) bool {
	for v := range q.C {
		if !pred(v) {
			return false
		}
	}
	return true
}

// First returns the first element matching pred, or zero+false.
func (q Query[T]) First(pred func(T) bool) (T, bool) {
	for v := range q.C {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Count returns the number of elements matching pred.
func (q Query[T]) Count(pred func(T) bool) int {
	cnt := 0
	for v := range q.C {
		if pred(v) {
			cnt++
		}
	}
	return cnt
}

// Len returns the total number of elements (identical to Count(func(T) bool { return true })).
func (q Query[T]) Len() int {
	cnt := 0
	for range q.C {
		cnt++
	}
	return cnt
}

// Reverse buffers all elements, then emits them in reverse order.
func (q Query[T]) Reverse() Query[T] {
	// buffer everything
	buf := q.ToSlice()
	return From(reverseSlice(buf))
}

// helper to reverse a slice
func reverseSlice[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
