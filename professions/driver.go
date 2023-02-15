package professions

import "github.com/eolinker/eosc"

var _ ITypedProfessionDrivers = (*TypedProfessionDrivers)(nil)

type ITypedProfessionDrivers interface {
	Get(name string) (eosc.IExtenderDriver, bool)
	Set(name string, d eosc.IExtenderDriver)
}
type TypedProfessionDrivers struct {
	eosc.Untyped[string, eosc.IExtenderDriver]
}

func NewProfessionDrivers() ITypedProfessionDrivers {
	return &TypedProfessionDrivers{Untyped: eosc.BuildUntyped[string, eosc.IExtenderDriver]()}
}
