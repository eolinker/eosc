package admin_o

type RollbackHandler interface {
	RollBack(api iAdminOperator) error
}
