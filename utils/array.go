package utils

func ArrayToMap[K comparable, V any](l []V, f func(V) K) map[K]V {
	r := make(map[K]V)
	for _, v := range l {
		r[f(v)] = v
	}
	return r
}
func GroupBy[K comparable, V any](l []V, f func(V) K) map[K][]V {
	r := make(map[K][]V)
	ArrayFilter(l, func(i int, v V) bool {
		k := f(v)
		if _, ok := r[k]; !ok {
			r[k] = make([]V, 0)
		}
		r[k] = append(r[k], v)
		return false
	})
	return r
}

func ArrayType[T any, V any](l []T, f func(T) V) []V {
	r := make([]V, len(l))
	for i, v := range l {
		r[i] = f(v)
	}
	return r
}
