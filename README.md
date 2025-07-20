# go-collection-boot
[![codecov](https://codecov.io/gh/SaiNageswarS/go-collection-boot/graph/badge.svg?token=XWI745R6EJ)](https://codecov.io/gh/SaiNageswarS/go-collection-boot)

**go-collection-boot** is a zero‑dependency, generics‑powered collection toolkit for Go 1.24+.
It brings ergonomic **LINQ‑style streaming**, **async helpers**, and rich **data‑structures**—all without reflection or external packages.

---

## ✨ Highlights

| Area              | Package | What you get                                                                                                          |
| ----------------- | ------- | --------------------------------------------------------------------------------------------------------------------- |
| Query & Transform | `linq`  | **High‑throughput streaming pipelines** with context‑aware cancellation (`Where`, `Select`, `Distinct`, `Flatten`, …) |
| Concurrency       | `async` | Type‑safe goroutine helpers (`Go`, `Await`, `AwaitAll`)                                                               |
| Collections       | `ds`    | Ergonomic generics for `Set`, `Stack`, `MinHeap`, and more                                                            |

---

## 🚀 Quick start

### LINQ – streaming pipelines that auto‑cancel

The `linq` package is built as a **heavy‑duty data‑flow engine**:

* **Streaming everywhere** – every transformer (`Where`, `Select`, `Flatten`, …) runs in its own goroutine, pushing items down an unbuffered/size‑hinted channel.
* **Compute‑heavy friendly** – downstream stages start before upstream finishes, so CPU‑bound work overlaps naturally.
* **Early termination** – sinks such as `First`, `Any`, or `Count` call `cancel()` as soon as their answer is known, short‑circuiting the entire chain and saving cycles.
* **Pure Go contexts** – the shared `context.Context` controls time‑outs, manual cancellation, and propagates errors upstream.

```go
package main

import (
    "context"
    "fmt"

    "github.com/SaiNageswarS/go-collection-boot/linq"
)

func main() {
    ctx := context.Background()
    nums := []int{1, 2, 3, 4, 5}

    // even squares ➜ slice
    squares, _ := linq.Pipe3(
        linq.FromSlice(ctx, nums),                 // Source  (Stream[int])
        linq.Where(func(n int) bool { return n%2==0 }),
        linq.Select(func(n int) int { return n*n }),
        linq.ToSlice[int](),                       // Sink  ([]int, error)
    )
    fmt.Println(squares) // [4 16]
}
```

#### Handy transformers & sinks

| Transformer       | Purpose                        |
| ----------------- | ------------------------------ |
| `Where(pred)`     | keep values matching predicate |
| `Select(mapFn)`   | map **T → U**                  |
| `Distinct(keyFn)` | deduplicate by key             |
| `Flatten[T]()`    | flatten `[][]T → []T`          |

| Sink                      | Returns         |
| ------------------------- | --------------- |
| `ToSlice[T]()`            | `([]T, error)`  |
| `Count[T]()`              | `(int, error)`  |
| `Any(pred)` / `All(pred)` | `(bool, error)` |
| `First[T]()`              | `(T, error)`    |
| `Reverse[T]()`            | `([]T, error)`  |

---

### Async – run tasks & await results 

```go
package main

import (
    "fmt"
    "github.com/SaiNageswarS/go-collection-boot/async"
)

func main() {
    fetchA := async.Go(func() (string, error) { return "alpha", nil })
    fetchB := async.Go(func() (string, error) { return "bravo", nil })

    words, err := async.AwaitAll(fetchA, fetchB)
    if err != nil { panic(err) }
    fmt.Println(words) // [alpha bravo]
}
```

---

### Data‑structures – Set, Stack, Min‑heap

```go
set := ds.NewSet[int]()
set.Add(1, 2, 3, 3)
fmt.Println(set.ToSlice()) // [1 2 3]

stk := ds.NewStack[string]()
stk.Push("first"); stk.Push("second")
val, _ := stk.Pop()
fmt.Println(val) // "second"

h := ds.NewMinHeap[int](func(a, b int) bool { return a < b })
h.Push(5); h.Push(2); h.Push(9)
fmt.Println(h.Pop()) // 2, true
```

---

## 🤝 Contributing

Pull requests are welcome! Run tests & `go vet ./...` before submitting.

---

## License

Apache‑2.0 © 2025 Sai Nageswar Satchidanand
