package slice

func Filter[T any](ts []T, f func(T) bool) []T {
	result := make([]T, 0)

	for _, v := range ts {
		if !f(v) {
			continue
		}

		result = append(result, v)
	}

	return result
}
