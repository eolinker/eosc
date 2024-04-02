package utils

func MapKey[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
func MapValue[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
func MapReGroup[K comparable, V any, G comparable](m map[K]V, f func(k K, v V) G) map[G][]V {
	r := make(map[G][]V)
	for k, v := range m {
		g := f(k, v)
		r[g] = append(r[g], v)
	}
	return r
}
func MapType[K comparable, V any, T any](m map[K]V, f func(k K, v V) (T, bool)) map[K]T {
	r := make(map[K]T)
	for k, v := range m {
		nv, yes := f(k, v)
		if yes {
			r[k] = nv
		}
	}
	return r
}
