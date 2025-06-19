package ds

type MinHeap[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewMinHeap[T any](less func(a, b T) bool) *MinHeap[T] {
	return &MinHeap[T]{
		data: []T{},
		less: less,
	}
}

func (h *MinHeap[T]) Push(value T) {
	h.data = append(h.data, value)
	h.up(len(h.data) - 1)
}

func (h *MinHeap[T]) Pop() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}

	min := h.data[0]
	last := len(h.data) - 1
	h.data[0] = h.data[last]
	h.data = h.data[:last]
	h.down(0)
	return min, true
}

func (h *MinHeap[T]) Peek() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}
	return h.data[0], true
}

func (h *MinHeap[T]) Len() int {
	return len(h.data)
}

func (h *MinHeap[T]) IsEmpty() bool {
	return len(h.data) == 0
}

func (h *MinHeap[T]) ToSlice() []T {
	return append([]T(nil), h.data...) // return a copy of the data slice
}

func (h *MinHeap[T]) up(idx int) {
	for idx > 0 {
		parent := (idx - 1) >> 1
		if !h.less(h.data[idx], h.data[parent]) {
			return
		}
		h.data[idx], h.data[parent] = h.data[parent], h.data[idx]
		idx = parent
	}
}

// heapify the subtree rooted at index
func (h *MinHeap[T]) down(idx int) {
	n := len(h.data)
	for {
		l, r := idx<<1+1, idx<<1+2
		smallest := idx
		if l < n && h.less(h.data[l], h.data[smallest]) {
			smallest = l
		}
		if r < n && h.less(h.data[r], h.data[smallest]) {
			smallest = r
		}
		if smallest == idx {
			return
		}
		h.data[idx], h.data[smallest] = h.data[smallest], h.data[idx]
		idx = smallest
	}
}
