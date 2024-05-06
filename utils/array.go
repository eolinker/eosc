package utils

import "github.com/eolinker/eosc/utils/set"

func ArrayToMap[K comparable, V any](l []V, f func(V) K) map[K]V {
	r := make(map[K]V)
	for _, v := range l {
		r[f(v)] = v
	}
	return r
}
func GroupBy[K comparable, V any](l []V, f func(V) K) map[K][]V {
	r := make(map[K][]V)
	for _, v := range l {
		k := f(v)
		if _, ok := r[k]; !ok {
			r[k] = make([]V, 0)
		}
		r[k] = append(r[k], v)
	}
	return r
}

func ArrayType[T any, V any](l []T, f func(T) V) []V {
	r := make([]V, len(l))
	for i, v := range l {
		r[i] = f(v)
	}
	return r
}

func Intersection[T comparable](a, b []T) []T {
	l := Min(len(a), len(b))
	if l == 0 {
		return nil
	}
	r := make([]T, 0, l)
	s := set.NewSet(a...)
	for _, bv := range b {
		if s.Contains(bv) {
			r = append(r, bv)
		}
	}
	return r
}
func Union[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	s := set.NewSet(a...)
	for _, bv := range b {
		s.Add(bv)
	}
	return s.List()
}
