package admin_o

type iAdminOperator interface {
	setWorker(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error)
	delWorker(id string) error
	setSetting(name string, data IData) error
	setVariable(namespace string, values map[string]string) error
}
