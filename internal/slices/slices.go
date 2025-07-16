package slices

import (
	"slices"
	"time"
)

func Transform[Slice ~[]S, S any, R any](source Slice, fn func(S) R) []R {
	result := make([]R, 0, len(source))
	for _, s := range source {
		result = append(result, fn(s))
	}
	return result
}

func SortByTimeDesc[Slice ~[]S, S any](source Slice, fn func(S) time.Time) {
	slices.SortFunc(source, func(a, b S) int {
		return fn(b).Compare(fn(a))
	})
}
