# go-collection-boot
[![codecov](https://codecov.io/gh/SaiNageswarS/go-collection-boot/graph/badge.svg?token=XWI745R6EJ)](https://codecov.io/gh/SaiNageswarS/go-collection-boot)

A zero-dependency, generics-powered collection library for Go.

This package offers utility types and functions to simplify working with collections — including async workflows, sets, and LINQ-style slice operations.

---

## ✨ Features

### ✅ `async` — Type-safe Goroutine Management

Run goroutines and await results without manual channels.

```go
import (
    "fmt"
    "github.com/SaiNageswarS/go-collection-boot/async"
)

func getEmbedding(ctx context.Context, text string) async.Result[[]float32] {
    return async.Go(func() ([]float32, error) {
        req := &EmbeddingRequest{
            Model: "text-embedding-3-small",
            Input: text,
        }

        // Assuming `client` is an initialized OpenAI client
        // resp has Embedding field of type []float32
        resp, err := client.CreateEmbedding(ctx, req) 
        if err != nil {
            return nil, err
        }

        return resp.Embedding, nil
    })
}

func main() {
    ctx := context.Background()
    text := "Go is awesome!"

    result, err := async.Await(getEmbedding(ctx, text))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Embedding:", result)
}
```

### ✅ `set` — Type-safe Set Operations
```go
import (
    "fmt"
    "github.com/SaiNageswarS/go-collection-boot/ds"
)

func main() {
    set := ds.NewSet[int]()
    set.Add(1, 2, 3)
    set.Add(3, 4) // Duplicate, will be ignored

    fmt.Println("Set contains 2:", set.Contains(2)) // true
    fmt.Println("Set size:", set.Len())            // 4

    // Iterate over elements
    for _, val := range set.ToSlice() {
		fmt.Println(val)
	}
}
```

### ✅ `linq` — LINQ-style Slice Operations
```go
import (
    "fmt"
    "github.com/SaiNageswarS/go-collection-boot/linq"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}

    // Filter even numbers and square them
    result := linq.From(numbers).
        Where(func(n int) bool { return n%2 == 0 }).
        Select(func(n int) int { return n * n }).
        ToSlice()

    fmt.Println("Squared even numbers:", result) // [4, 16]
}
``` 

## Installation

```bash
go get github.com/SaiNageswarS/go-collection-boot
```