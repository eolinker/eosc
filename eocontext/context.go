package eocontext

import (
	"context"
)

type CompleteHandler interface {
	Complete(ctx EoContext) error
}

type FinishHandler interface {
	Finish(ctx EoContext) error
}

type EoContext interface {
	RequestId() string
	Context() context.Context
	Value(key interface{}) interface{}
	WithValue(key, val interface{})

	Scheme() string
	Assert(i interface{}) error

	SetLabel(name, value string)
	GetLabel(name string) string
	Labels() map[string]string

	GetComplete() CompleteHandler
	SetCompleteHandler(handler CompleteHandler)
	GetFinish() FinishHandler
	SetFinish(handler FinishHandler)
	GetApp() EoApp
	SetApp(app EoApp)
	GetBalance() BalanceHandler
	SetBalance(handler BalanceHandler)
}
