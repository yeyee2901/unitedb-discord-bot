package helper

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// RemoveDuplicate is used to remove any duplicate element in the entry
func RemoveDuplicate[T constraints.Ordered](in []T) []T {
	entry := make(map[T]bool)
	out := []T{}

	for idx := range in {
		if _, exist := entry[in[idx]]; !exist {
			entry[in[idx]] = true
			out = append(out, in[idx])
		}
	}

	return out
}

// SortUsingOrder is used to sort the input array according to the order array
// priority
func SortUsingOrder[T constraints.Ordered](in []T, order []T) ([]T, error) {
	out := []T{}

	for numOrder := range order {
		for arrIdx := range in {
			if in[arrIdx] == order[numOrder] {
				out = append(out, in[arrIdx])
			}
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("Empty sorted result.")
	}

	return out, nil
}
