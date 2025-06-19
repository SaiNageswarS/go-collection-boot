package linq

// These functions provide utility methods for working with slices of any type.

func Flatten[T any](items [][]T) []T {
	var result []T
	for _, sublist := range items {
		result = append(result, sublist...)
	}
	return result
}

func Distinct[T any, K comparable](items []T, keySelector func(T) K) []T {
	seen := make(map[K]struct{})
	var result []T
	for _, item := range items {
		key := keySelector(item)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}
	return result
}
