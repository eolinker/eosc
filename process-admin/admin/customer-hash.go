package admin

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils/hash"
)

var (
	_ eosc.ICustomerVar = (*imlCustomerHash)(nil)
)

type stringHash = hash.Hash[string, string]

type imlCustomerHash struct {
	data eosc.Untyped[string, stringHash]
}

func newImlCustomerHash(data eosc.Untyped[string, stringHash]) *imlCustomerHash {
	return &imlCustomerHash{data: data}
}

func (i *imlCustomerHash) Get(key string, field string) (string, bool) {
	vs, has := i.data.Get(key)
	if !has {
		return "", false
	}
	return vs.Get(field)
}

func (i *imlCustomerHash) GetAll(key string) (map[string]string, bool) {
	vs, has := i.data.Get(key)
	if !has {
		return nil, false
	}
	return vs.Map(), true
}

func (i *imlCustomerHash) Exists(key string, field string) bool {
	_, has := i.Get(key, field)
	return has
}
