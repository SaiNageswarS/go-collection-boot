package linq

import (
	"context"
	"errors"
)

// Stream is a lazy, cancel‑aware sequence of T.
type Stream[T any] struct {
	ctx    context.Context    // shared cancellation context
	cancel context.CancelFunc // cancel function for short-circuit operations like All/Any/First
	cap    int                // capacity hint for the channel. Advisory only, downstream can shrink it.
	C      <-chan T           // receive‑only element channel
}

func NewStream[T any](ctx context.Context, out <-chan T, cancel context.CancelFunc, capHint ...int) Stream[T] {
	if len(capHint) > 0 && capHint[0] > 0 {
		return Stream[T]{ctx: ctx, C: out, cancel: cancel, cap: capHint[0]}
	}

	return Stream[T]{ctx: ctx, C: out, cancel: cancel, cap: 1}
}

// ---------- 1 · Sources ----------

// FromSlice starts a new stream that emits the items one by one.
func FromSlice[T any](parent context.Context, src []T) Stream[T] {
	ctx, cancel := context.WithCancel(parent)

	out := make(chan T, max(1, len(src)/2))
	go func() {
		defer close(out)
		for _, v := range src {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()

	return NewStream(ctx, out, cancel, len(src))
}

func trySend[T any](ctx context.Context, out chan<- T, v T) bool {
	select {
	case <-ctx.Done():
		return false // downstream cancelled
	case out <- v:
		return true // delivered
	}
}

func tryRecv[T any](ctx context.Context, in <-chan T) (v T, ok bool, err error) {
	select {
	case <-ctx.Done():
		return v, false, ctx.Err()
	case v, ok = <-in:
		return v, ok, nil
	}
}

// ---------- 2 · Transformers ----------

// Where keeps only the values for which pred == true.
func whereFn[T any](in Stream[T], pred func(T) bool) Stream[T] {
	out := make(chan T, max(1, in.cap/2))

	go func() {
		defer close(out)
		for {
			v, ok, err := tryRecv(in.ctx, in.C)
			if err != nil { // context cancelled
				return
			}
			if !ok { // channel closed
				return
			}

			if pred(v) {
				if !trySend(in.ctx, out, v) {
					return // downstream cancelled
				}
			}
		}
	}()

	return NewStream(in.ctx, out, in.cancel, in.cap)
}

// Select maps T → U.
func selectFn[T, U any](in Stream[T], f func(T) U) Stream[U] {
	out := make(chan U, max(1, in.cap/2))

	go func() {
		defer close(out)
		for {
			v, ok, err := tryRecv(in.ctx, in.C)
			if err != nil { // context cancelled
				return
			}
			if !ok { // channel closed
				return
			}

			if !trySend(in.ctx, out, f(v)) {
				return
			}
		}
	}()

	return NewStream(in.ctx, out, in.cancel, in.cap)
}

func distinctFn[T any, K comparable](in Stream[T], keySelector func(T) K) Stream[T] {
	out := make(chan T, max(1, in.cap/2))
	seen := make(map[K]struct{})

	go func() {
		defer close(out)
		for {
			v, ok, err := tryRecv(in.ctx, in.C)
			if err != nil { // context cancelled
				return
			}
			if !ok { // channel closed
				return
			}

			key := keySelector(v)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				if !trySend(in.ctx, out, v) {
					return // downstream cancelled
				}
			}
		}
	}()

	return NewStream(in.ctx, out, in.cancel, in.cap)
}

func flattenFn[T any](in Stream[[]T]) Stream[T] {
	out := make(chan T, 64) // heuristic buffer size since we don't know the size of the inner slices

	go func() {
		defer close(out)
		for {
			v, ok, err := tryRecv(in.ctx, in.C)
			if err != nil { // context cancelled
				return
			}
			if !ok { // channel closed
				return
			}

			for _, item := range v {
				if !trySend(in.ctx, out, item) {
					return // downstream cancelled
				}
			}
		}
	}()

	return NewStream(in.ctx, out, in.cancel, in.cap)
}

// ---------- 3 · Sinks ----------

// ToSlice collects every remaining element (honours ctx cancellation).
func toSliceFn[T any](s Stream[T]) ([]T, error) {
	defer s.cancel()

	out := make([]T, 0, s.cap)

	for {
		v, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return out, err
		}
		if !ok { // channel closed
			return out, nil
		}
		out = append(out, v)
	}
}

// ForEach applies fn to every element; stops early if ctx is cancelled.
func forEachFn[T any](s Stream[T], fn func(T)) error {
	defer s.cancel()

	for {
		v, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return err
		}
		if !ok { // channel closed
			return nil
		}
		fn(v)
	}
}

func countFn[T any](s Stream[T]) (int, error) {
	defer s.cancel()

	count := 0
	for {
		_, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return count, err
		}
		if !ok { // channel closed
			return count, nil
		}
		count++
	}
}

func firstFn[T any](s Stream[T]) (T, error) {
	defer s.cancel() // close upstream transformers

	var zero T
	v, ok, err := tryRecv(s.ctx, s.C)
	if err != nil { // context cancelled
		return zero, err
	}
	if !ok { // channel closed
		return zero, errors.New("stream is empty, no first element found")
	}

	return v, nil // return the first element found
}

func reverseFn[T any](s Stream[T]) ([]T, error) {
	defer s.cancel() // close upstream transformers

	out := make([]T, 0, s.cap)

	for {
		v, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return nil, err
		}
		if !ok { // channel closed
			break
		}
		out = append(out, v)
	}

	// Reverse the slice in place.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out, nil
}

func allFn[T any](s Stream[T], pred func(T) bool) (bool, error) {
	defer s.cancel() // close upstream transformers

	for {
		v, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return false, err
		}
		if !ok { // channel closed
			return true, nil
		}
		if !pred(v) {
			return false, nil // found a value that does not match the predicate
		}
	}
}

func anyFn[T any](s Stream[T], pred func(T) bool) (bool, error) {
	defer s.cancel() // close upstream transformers

	for {
		v, ok, err := tryRecv(s.ctx, s.C)
		if err != nil { // context cancelled
			return false, err
		}
		if !ok { // channel closed
			return false, nil // no matching value found
		}
		if pred(v) {
			return true, nil // found a value that matches the predicate
		}
	}
}

// ---------- 4 · Public "curried" Adapters ----------

func Where[T any](pred func(T) bool) func(Stream[T]) Stream[T] {
	return func(in Stream[T]) Stream[T] {
		return whereFn(in, pred)
	}
}

func Select[T, U any](f func(T) U) func(Stream[T]) Stream[U] {
	return func(in Stream[T]) Stream[U] {
		return selectFn(in, f)
	}
}

func Distinct[T any, K comparable](keySelector func(T) K) func(Stream[T]) Stream[T] {
	return func(in Stream[T]) Stream[T] {
		return distinctFn(in, keySelector)
	}
}

func Flatten[T any]() func(Stream[[]T]) Stream[T] {
	return func(in Stream[[]T]) Stream[T] {
		return flattenFn(in)
	}
}

func ToSlice[T any]() func(Stream[T]) ([]T, error) {
	return func(s Stream[T]) ([]T, error) {
		return toSliceFn(s)
	}
}

func ForEach[T any](fn func(T)) func(Stream[T]) error {
	return func(s Stream[T]) error {
		return forEachFn(s, fn)
	}
}

func Count[T any]() func(Stream[T]) (int, error) {
	return func(s Stream[T]) (int, error) {
		return countFn(s)
	}
}

func First[T any]() func(Stream[T]) (T, error) {
	return func(s Stream[T]) (T, error) {
		return firstFn(s)
	}
}

func Reverse[T any]() func(Stream[T]) ([]T, error) {
	return func(s Stream[T]) ([]T, error) {
		return reverseFn(s)
	}
}

func All[T any](pred func(T) bool) func(Stream[T]) (bool, error) {
	return func(s Stream[T]) (bool, error) {
		return allFn(s, pred)
	}
}

func Any[T any](pred func(T) bool) func(Stream[T]) (bool, error) {
	return func(s Stream[T]) (bool, error) {
		return anyFn(s, pred)
	}
}
