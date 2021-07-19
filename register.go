package eosc

import (
	"github.com/eolinker/eosc/internal"
)

type IRegister interface {
	Register(name string, obj interface{}, force bool) error
	Get(name string) (interface{}, bool)
}

type IRegisterData interface {
	Set(name string, v interface{})
	Get(name string) (interface{}, bool)
}

type Register struct {
	data IRegisterData
}

func NewRegister() IRegister {
	return &Register{
		data: internal.NewUntyped(),
	}
}

func (r *Register) Register(name string, obj interface{}, force bool) error {

	if _, has := r.data.Get(name); has && !force {
		return ErrorRegisterConflict
	}
	r.data.Set(name, obj)
	return nil
}

func (r *Register) Get(name string) (interface{}, bool) {
	return r.data.Get(name)
}
