package eosc

type IRegister interface {
	Register(name string, obj interface{}, force bool) error
	Get(name string) (interface{}, bool)
	Del(name string) (interface{}, bool)
}

type Register struct {
	data IUntyped
}

func (r *Register) Del(name string) (interface{}, bool) {
	return r.data.Del(name)
}

func NewRegister() IRegister {
	return &Register{
		data: NewUntyped(),
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
