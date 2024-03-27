package admin_o

type rollBackForCreate string

func (r rollBackForCreate) RollBack(data iAdminOperator) error {

	return data.delWorker(string(r))
}

type rollBackForSet struct {
}
