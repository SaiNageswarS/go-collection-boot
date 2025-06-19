package linq

type Query[T any] struct {
	items []T
}

func From[T any](items []T) Query[T] {
	return Query[T]{items: items}
}

func (q Query[T]) Where(predicate func(T) bool) Query[T] {
	var result []T
	for _, item := range q.items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return Query[T]{items: result}
}

func (q Query[T]) Select(mapper func(T) T) Query[T] {
	var result []T
	for _, item := range q.items {
		result = append(result, mapper(item))
	}
	return Query[T]{items: result}
}

func (q Query[T]) Any(predicate func(T) bool) bool {
	for _, item := range q.items {
		if predicate(item) {
			return true
		}
	}
	return false
}

func (q Query[T]) All(predicate func(T) bool) bool {
	for _, item := range q.items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func (q Query[T]) First(predicate func(T) bool) (T, bool) {
	for _, item := range q.items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T // zero value for type T
	return zero, false
}

func (q Query[T]) Count(predicate func(T) bool) int {
	count := 0
	for _, item := range q.items {
		if predicate(item) {
			count++
		}
	}
	return count
}

func (q Query[T]) Len() int {
	return len(q.items)
}

func (q Query[T]) Reverse() Query[T] {
	var result []T
	for i := len(q.items) - 1; i >= 0; i-- {
		result = append(result, q.items[i])
	}
	return Query[T]{items: result}
}

// ToSlice converts the Query to a slice.
func (q Query[T]) ToSlice() []T {
	return q.items
}
