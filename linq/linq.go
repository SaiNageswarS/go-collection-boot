package linq

// Where filters a slice based on a predicate.
func Where[T any](input []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range input {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Select maps each item in a slice to a new form.
func Select[T any, R any](input []T, mapper func(T) R) []R {
	result := make([]R, 0, len(input))
	for _, item := range input {
		result = append(result, mapper(item))
	}
	return result
}

// Any returns true if any element matches the predicate.
func Any[T any](input []T, predicate func(T) bool) bool {
	for _, item := range input {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All returns true if all elements match the predicate.
func All[T any](input []T, predicate func(T) bool) bool {
	for _, item := range input {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// First returns the first element that matches the predicate.
func First[T any](input []T, predicate func(T) bool) (T, bool) {
	for _, item := range input {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Count returns number of elements matching a predicate.
func Count[T any](input []T, predicate func(T) bool) int {
	count := 0
	for _, item := range input {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Distinct returns a slice with duplicate values removed.
func Distinct[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	var result []T
	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Reverse returns a new slice with the elements reversed.
func Reverse[T any](input []T) []T {
	n := len(input)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = input[n-1-i]
	}
	return result
}
