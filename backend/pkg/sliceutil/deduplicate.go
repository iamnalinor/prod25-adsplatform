package sliceutil

import "slices"

// DeduplicateLast removes duplicates from slice, keeping only the last ones.
func DeduplicateLast[S ~[]T, T comparable](slice S) S {
	res := make(S, 0, len(slice))
	exists := make(map[T]bool)
	for i := len(slice) - 1; i >= 0; i-- {
		if !exists[slice[i]] {
			res = append(res, slice[i])
			exists[slice[i]] = true
		}
	}
	slices.Reverse(res)
	return res
}
