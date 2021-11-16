package extends

import (
	"fmt"

	"github.com/eolinker/eosc"
)

type ExtenderRegister struct {
	group   string
	project string
	data    map[string]eosc.IExtenderDriverFactory
}

func NewExtenderRegister(group string, project string) *ExtenderRegister {
	return &ExtenderRegister{group: group, project: project, data: make(map[string]eosc.IExtenderDriverFactory)}
}

func (r *ExtenderRegister) RegisterExtenderDriver(name string, factory eosc.IExtenderDriverFactory) error {
	_, has := r.data[name]
	if has {
		return fmt.Errorf("%s:%w", name, ErrorExtenderNameDuplicate)
	}
	r.data[name] = factory
	return nil
}

func (r *ExtenderRegister) RegisterTo(register eosc.IExtenderDriverRegister) {
	for n, f := range r.data {
		id := FormatDriverId(r.group, r.project, n)
		register.RegisterExtenderDriver(id, f)
	}
}

func (r *ExtenderRegister) All() []string {
	rs := make([]string, 0, len(r.data))
	for n := range r.data {
		rs = append(rs, n)
	}
	return rs
}
