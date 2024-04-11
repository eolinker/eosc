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
	_, err := api.setWorker((*eosc.WorkerConfig)(r))
	if err != nil {
		return err
	}
	return nil
}

type rollbackForDelete WorkerInfo

func (r *rollbackForDelete) RollBack(api iAdminOperator) error {
	_, err := api.setWorker(r.config)
	if err != nil {
		return err
	}
	return nil
}

func newRollbackForDelete(worker *WorkerInfo) RollbackHandler {
	return (*rollbackForDelete)(worker)
}
