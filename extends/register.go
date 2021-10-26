package extends

import (
	"fmt"

	"github.com/eolinker/eosc"
)

type ExtenderRegister struct {
	group string
	name  string
	data  map[string]eosc.IExtenderDriverFactory
}

func NewExtenderRegister(group string, name string) *ExtenderRegister {
	return &ExtenderRegister{group: group, name: name, data: make(map[string]eosc.IExtenderDriverFactory)}
}

func (r *ExtenderRegister) RegisterExtender(name string, factory eosc.IExtenderDriverFactory) error {
	_, has := r.data[name]
	if has {
		return fmt.Errorf("%s:%w", name, ErrorExtenderNameDuplicate)
	}
	r.data[name] = factory
	return nil
}

func (r *ExtenderRegister) RegisterTo(register eosc.IExtenderRegister) {
	for n, f := range r.data {
		id := FormatDriverId(r.group, r.name, n)
		register.RegisterExtender(id, f)
	}
}

func (r *ExtenderRegister) All() []string {
	rs := make([]string, 0, len(r.data))
	for n := range r.data {
		rs = append(rs, FormatDriverId(r.group, r.name, n))
	}
	return rs
}
