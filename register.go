package eosc

type IRegister[T any] interface {
	Register(name string, obj T, force bool) error
	Get(name string) (T, bool)
	Del(name string) (T, bool)
}
type Register[T any] struct {
	Untyped[string, T]
}

func NewRegister[T any]() IRegister[T] {
	return &Register[T]{
		Untyped: BuildUntyped[string, T](),
	}
}

func (r *Register[T]) Register(name string, obj T, force bool) error {

	if _, has := r.Untyped.Get(name); has && !force {
		return ErrorRegisterConflict
	}
	r.Untyped.Set(name, obj)
	return nil
}
