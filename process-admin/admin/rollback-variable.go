package admin

type rollBackForCreateVariable string

func (r rollBackForCreateVariable) RollBack(api iAdminOperator) error {
	return api.setVariable(string(r), nil)
}

func newRollbackForCreatVariable(namespace string) RollbackHandler {
	return rollBackForCreateVariable(namespace)
}

type rollbackForUpdateVariable struct {
	namespace string
	value     map[string]string
}

func (r *rollbackForUpdateVariable) RollBack(api iAdminOperator) error {
	return api.setVariable(r.namespace, r.value)
}

func newRollbackForUpdateVariable(namespace string, value map[string]string) RollbackHandler {
	return &rollbackForUpdateVariable{namespace: namespace, value: value}
}
