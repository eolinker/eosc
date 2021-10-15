package process_worker

import "github.com/eolinker/eosc"

var _ ITypedDrivers = (*TypedDrivers)(nil)

type ITypedDrivers interface {
	Get(name string) (eosc.IProfessionDriver, bool)
}
type TypedDrivers struct {
	data eosc.IUntyped
}

func NewTypedDrivers() *TypedDrivers {
	return &TypedDrivers{data: eosc.NewUntyped()}
}

func (t *TypedDrivers) Get(name string) (eosc.IProfessionDriver, bool) {
	if o, has := t.data.Get(name); has {
		return o.(eosc.IProfessionDriver), has
	}
	return nil, false
}
