package async

import (
	"fmt"
	"runtime/debug"
)

type Result[T any] struct {
	Data T
	Err  error
}

func Await[T any](ch <-chan Result[T]) (T, error) {
	res := <-ch
	return res.Data, res.Err
}

func AwaitAll[T any](chs ...<-chan Result[T]) ([]T, error) {
	type pair struct {
		i int
		r Result[T]
	}

	out := make(chan pair, len(chs))
	for i, ch := range chs {
		i, ch := i, ch // capture copy
		go func() { out <- pair{i, <-ch} }()
	}

	results := make([]T, len(chs))
	var firstErr error
	for range chs {
		p := <-out // drains every channel
		if firstErr == nil {
			firstErr = p.r.Err // remember first error
		}
		results[p.i] = p.r.Data
	}
	return results, firstErr
}

// PanicError is the error Go delivers when fn panics.
type PanicError struct {
	Value any    // the value passed to panic
	Stack []byte // the goroutine's stack at the panic
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("async: panic: %v\n%s", e.Value, e.Stack)
}

// Unwrap returns the panic value when it is an error, so errors.Is and
// errors.As see through the panic.
func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error)
	return err
}

// Go runs fn in a new goroutine and delivers its result on the returned
// channel. A panic in fn is recovered and delivered as a *PanicError instead
// of crashing the process.
func Go[T any](fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1) // buffered so sender never blocks
	go func() {
		defer close(ch)
		ch <- run(fn)
	}()
	return ch
}

func run[T any](fn func() (T, error)) (res Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			res = Result[T]{Err: &PanicError{Value: r, Stack: debug.Stack()}}
		}
	}()
	v, err := fn()
	return Result[T]{Data: v, Err: err}
}
