package eosc

type ExtenderBuilder interface {
	Register(register IExtenderDriverRegister)
}
