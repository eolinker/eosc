package admin

import "github.com/eolinker/eosc"

type roleBackForAddProfession string

func newRollbackForAddProfession(name string) RollbackHandler {
	return roleBackForAddProfession(name)
}
func (r roleBackForAddProfession) RollBack(api iAdminOperator) error {
	return api.delProfession(string(r))
}

type rollbackForSetProfession struct {
	name   string
	config *eosc.ProfessionConfig
}

func newRollbackForSetProfession(name string, config *eosc.ProfessionConfig) *rollbackForSetProfession {
	return &rollbackForSetProfession{name: name, config: config}
}

func (r *rollbackForSetProfession) RollBack(api iAdminOperator) error {
	return api.setProfession(r.name, r.config)
}
