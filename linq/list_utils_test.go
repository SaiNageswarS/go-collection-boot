package linq

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	data := [][]int{{1, 2}, {3, 4}, {5}}
	flattened := Flatten(data)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, flattened, "Flatten should concatenate sublists into a single list")
}

func TestDistinct(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 3}
	distinct := Distinct(data, func(n int) int { return n })
	expected := []int{1, 2, 3}
	assert.Equal(t, expected, distinct, "Distinct should remove duplicates based on the key selector")
}

func TestDistinctWithKeySelector(t *testing.T) {
	data := []struct {
		ID   int
		Name string
	}{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice"}, // duplicate by ID
	}

	distinct := Distinct(data, func(item struct {
		ID   int
		Name string
	}) int {
		return item.ID
	})

	expected := []struct {
		ID   int
		Name string
	}{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	assert.Equal(t, expected, distinct, "Distinct with key selector should remove duplicates based on the key")
}

func TestMap(t *testing.T) {
	data := []int{1, 2, 3}
	mapped := Map(data, func(n int) string { return "id_" + strconv.Itoa(n) })
	expected := []string{"id_1", "id_2", "id_3"}
	assert.Equal(t, expected, mapped, "Map should transform each item in the list")
}
