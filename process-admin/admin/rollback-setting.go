package admin

type rollbackForSettingSet struct {
	config []byte
	name   string
}

func newRollbackForSettingSet(name string, config []byte) RollbackHandler {
	return &rollbackForSettingSet{config: config, name: name}
}

func (r *rollbackForSettingSet) RollBack(api iAdminOperator) error {
	return api.setSetting(r.name, r.config)
}
