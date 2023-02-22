package eosc

import (
	"fmt"
)

var (
// DefaultProfessionDriverRegister IExtenderDriverRegister = NewExtenderRegister()
)

type ExtenderRegister struct {
	data IRegister[IExtenderDriverFactory]
}

func (p *ExtenderRegister) RegisterExtenderDriver(name string, factory IExtenderDriverFactory) error {

	err := p.data.Register(name, factory, false)
	if err != nil {
		return fmt.Errorf("register profession  driver %s:%w", name, err)
	}
	return nil
}

func (p *ExtenderRegister) GetDriver(name string) (IExtenderDriverFactory, bool) {

	if v, has := p.data.Get(name); has {
		return v, true
	}
	return nil, false
}

func NewExtenderRegister() *ExtenderRegister {
	return &ExtenderRegister{
		data: NewRegister[IExtenderDriverFactory](),
	}
}

type IExtenderDriverRegister interface {
	RegisterExtenderDriver(name string, factory IExtenderDriverFactory) error
}

type IExtenderDrivers interface {
	GetDriver(name string) (IExtenderDriverFactory, bool)
}
type IExtenderDriverManager interface {
	IExtenderDriverRegister
}
