package eocontext

import "time"

type BalanceHandler interface {
	Select(ctx EoContext) (INode, int, error)
	Scheme() string
	TimeOut() time.Duration
	EoApp
}
