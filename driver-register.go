package eosc

import "fmt"
var(
	DefaultProfessionDriverRegister IDriverRegister = NewProfessionDriverRegister()
)

type ProfessionDriverRegister struct {
	data IRegister
}

func (p *ProfessionDriverRegister) RegisterProfessionDriver(name string, factory IProfessionDriverFactory) error {
	err := p.data.Register(name, factory, false)
	if err!=nil{
		return fmt.Errorf("register profession  driver %s:%w",name,err)
	}
	return  nil
}

func (p *ProfessionDriverRegister) GetProfessionDriver(name string) (IProfessionDriverFactory, bool) {
	if v, has := p.data.Get(name);has{
		return v.(IProfessionDriverFactory),true
	}
	return nil,false
}

func NewProfessionDriverRegister() *ProfessionDriverRegister {
	return &ProfessionDriverRegister{
		data:NewRegister(),
	}
}

type IDriverRegister interface {
	RegisterProfessionDriver(name string,factory IProfessionDriverFactory) error
	GetProfessionDriver(name string)(IProfessionDriverFactory,bool)
}

func RegisterProfessionDriver(name string, factory IProfessionDriverFactory) error {

	return  DefaultProfessionDriverRegister.RegisterProfessionDriver(name, factory)
}

func GetProfessionDriver(name string) (IProfessionDriverFactory, bool) {

	return  DefaultProfessionDriverRegister.GetProfessionDriver(name)
}
