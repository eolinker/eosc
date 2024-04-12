/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package set

var (
	_ Set[int] = (*imlSet[int])(nil)
)

type Set[T comparable] interface {
	Contains(i T) bool
	Add(i T)
	Remove(i T)
	Size() int
	Clear()
	List() []T
}
type imlSet[T comparable] struct {
	data map[T]struct{}
}

func (s *imlSet[T]) List() []T {
	r := make([]T, 0, len(s.data))
	for k := range s.data {
		r = append(r, k)
	}
	return r
}

func (s *imlSet[T]) Contains(i T) bool {
	_, has := s.data[i]
	return has
}

func (s *imlSet[T]) Add(i T) {
	s.data[i] = struct{}{}
}

func (s *imlSet[T]) Remove(i T) {
	delete(s.data, i)
}

func (s *imlSet[T]) Size() int {
	return len(s.data)
}

func (s *imlSet[T]) Clear() {
	s.data = make(map[T]struct{})
}

func NewSet[T comparable](is ...T) Set[T] {

	s := &imlSet[T]{data: make(map[T]struct{})}
	for _, i := range is {
		s.Add(i)
	}
	return s
}
func NewMapSet[T comparable, V any](m map[T]V) Set[T] {

	s := &imlSet[T]{data: make(map[T]struct{})}
	for k, _ := range m {
		s.Add(k)
	}
	return s
}
