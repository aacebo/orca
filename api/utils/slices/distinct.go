package slices

func Distinct[T comparable](v []T) []T {
	m := map[T]bool{}
	arr := []T{}

	for _, item := range v {
		if _, exists := m[item]; exists {
			continue
		}

		arr = append(arr, item)
		m[item] = true
	}

	return arr
}
