package eosc

import (
	"plugin"

	"github.com/eolinker/eosc/log"
)

type RegisterFunc func(register IExtenderRegister)

func LoadExtender(file string, register IExtenderRegister) error {

	p, err := plugin.Open(file)
	if err != nil {
		log.Errorf("error to open plugin %s:%s", file, err.Error())
		return err
	}

	r, err := p.Lookup("Register")
	if err != nil {
		log.Errorf("call register from  plugin : %s : %s", file, err.Error())
		return err
	}
	r.(RegisterFunc)(register)
	return nil
}
