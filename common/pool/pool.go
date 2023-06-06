/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package pool

import "sync"

type Pool[T any] interface {
	Get() T
	PUT(t T)
}

type _Pool[T any] struct {
	pool sync.Pool
}

func (p *_Pool[T]) Get() T {
	v := p.pool.Get()
	return v.(T)
}

func (p *_Pool[T]) PUT(t T) {
	p.pool.Put(t)
}
func New[T any](new func() T) Pool[T] {
	return &_Pool[T]{
		pool: sync.Pool{New: func() any { return new() }},
	}
}
