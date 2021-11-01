package process_worker

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/extends"
)

func loadPluginEnv() eosc.IExtenderDrivers {
	register := eosc.NewExtenderRegister()
	extends.LoadInner(register)
	return register
}
