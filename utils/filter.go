package utils

func ArrayFilter[T any](a []T, filter func(i int, v T) bool) []T {
	r := make([]T, 0, len(a))
	for i, v := range a {
		if filter(i, v) {
			r = append(r, v)
		}
	}
	return r
}
func MapFilter[K comparable, V any](m map[K]V, filter func(k K, v V) bool) map[K]V {
	r := make(map[K]V, len(m))
	for k, v := range m {
		if filter(k, v) {
			r[k] = v
		}
	}
	return r
}
