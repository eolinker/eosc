package professions

import "github.com/eolinker/eosc"

var _ ITypedProfessionDrivers = (*TypedProfessionDrivers)(nil)

type ITypedProfessionDrivers interface {
	Get(name string) (eosc.IExtenderDriver, bool)
}
type TypedProfessionDrivers struct {
	data eosc.IUntyped
}

func NewProfessionDrivers() *TypedProfessionDrivers {
	return &TypedProfessionDrivers{data: eosc.NewUntyped()}
}

func (t *TypedProfessionDrivers) Get(name string) (eosc.IExtenderDriver, bool) {
	if o, has := t.data.Get(name); has {
		return o.(eosc.IExtenderDriver), has
	}
	return nil, false
}
