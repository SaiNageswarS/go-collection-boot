package ds

type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet creates a new empty Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

// FromSlice creates a new Set from a slice.
func FromSlice[T comparable](items []T) *Set[T] {
	s := NewSet[T]()
	s.AddAll(items)
	return s
}

// Add inserts one or more items into the set.
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

// AddAll adds all elements from a slice into the set.
func (s *Set[T]) AddAll(items []T) {
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

// Remove deletes an item from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}

// Contains checks whether the item exists in the set.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

// ContainsAll returns true if all items are in the set.
func (s *Set[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ToSlice returns all elements as a slice.
func (s *Set[T]) ToSlice() []T {
	out := make([]T, 0, len(s.data))
	for item := range s.data {
		out = append(out, item)
	}
	return out
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.data)
}

// Clear removes all elements.
func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
}

// IsEmpty checks if the set is empty.
func (s *Set[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// Clone returns a copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	clone := NewSet[T]()
	for item := range s.data {
		clone.Add(item)
	}
	return clone
}

// Union returns a new set that is the union of two sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	out := s.Clone()
	for item := range other.data {
		out.Add(item)
	}
	return out
}

// Intersection returns a new set with common elements.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	out := NewSet[T]()
	for item := range s.data {
		if other.Contains(item) {
			out.Add(item)
		}
	}
	return out
}

// Difference returns a new set with items in s but not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	out := NewSet[T]()
	for item := range s.data {
		if !other.Contains(item) {
			out.Add(item)
		}
	}
	return out
}
