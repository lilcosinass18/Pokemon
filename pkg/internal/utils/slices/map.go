package slices

func Map[X, Y any](arr []X, f func(X) Y) []Y {
	result := make([]Y, 0, len(arr))

	for _, elem := range arr {
		result = append(result, f(elem))
	}

	return result
}
