package slice

import slices2 "slices"

func Unique[T comparable](ts []T) []T {
	return UniqueFunc(ts, func(a T, b T) bool {
		return a == b
	})
}

func UniqueFunc[T any](ts []T, fn func(a T, b T) bool) []T {
	uTs := make([]T, 0)

	for _, t1 := range ts {
		if slices2.ContainsFunc(uTs, func(t2 T) bool {
			return fn(t1, t2)
		}) {
			continue
		}

		uTs = append(uTs, t1)
	}

	return uTs
}
