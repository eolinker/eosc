package admin

import "github.com/eolinker/eosc"

type rollBackForCreate string

func newRollBackForCreate(id string) rollBackForCreate {
	return rollBackForCreate(id)
}
func (r rollBackForCreate) RollBack(data iAdminOperator) error {

	_, err := data.delWorker(string(r))
	if err != nil {
		return err
	}
	return nil
}

type rollBackForSet eosc.WorkerConfig

func newRollBackForSet(worker *eosc.WorkerConfig) RollbackHandler {
	return (*rollBackForSet)(worker)
}

func (r *rollBackForSet) RollBack(api iAdminOperator) error {
	_, err := api.setWorker(r.Id, r.Profession, r.Name, r.Driver, r.Version, r.Description, r.Body, r.Update, r.Create)
	if err != nil {
		return err
	}
	return nil
}

type rollbackForDelete WorkerInfo

func (r *rollbackForDelete) RollBack(api iAdminOperator) error {
	_, err := api.setWorker(r.config.Id, r.config.Profession, r.config.Name, r.config.Driver, r.config.Version, r.config.Description, r.config.Body, r.config.Update, r.config.Create)
	if err != nil {
		return err
	}
	return nil
}

func newRollbackForDelete(worker *WorkerInfo) RollbackHandler {
	return (*rollbackForDelete)(worker)
}
