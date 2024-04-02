package admin

type RollbackHandler interface {
	RollBack(api iAdminOperator) error
}

type RollbackHandlerList []RollbackHandler

func (r RollbackHandlerList) RollBack(api iAdminOperator) error {
	for i := len(r) - 1; i >= 0; i-- {
		if err := r[i].RollBack(api); err != nil {
			return err
		}
	}
	return nil
}
