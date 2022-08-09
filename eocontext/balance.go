package eocontext

type BalanceHandler interface {
	Select(ctx EoContext) (INode, error)
}
