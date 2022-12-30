package helper

import "fmt"

// generic function to remove duplicate element in an array that is
// comparable (string, int, bool, etc)
func RemoveDuplicateElement[T comparable](arr []T) (result []T) {
	entry := make(map[T]bool)

	for index := range arr {
		if _, exist := entry[arr[index]]; !exist {
			entry[arr[index]] = true
			result = append(result, arr[index])
		}
	}

	return
}

// Generic function to sort an array based on the order given
// NOTE: If the result is empty, then it's considered failed sorting.
// NOTE: This can be caused by no elements matched the ordering array
func SortUsingOrderRule[T comparable](arr []T, orderRule []T) (result []T, err error) {
	for numOrder := range orderRule {
		for arrIdx := range arr {
			if arr[arrIdx] == orderRule[numOrder] {
				result = append(result, arr[arrIdx])
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Empty sorted result.")
	}

	return result, nil
}
