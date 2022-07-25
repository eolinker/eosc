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
	Complete() CompleteHandler
	SetCompleteHandler(handler CompleteHandler)
	Finish() FinishHandler
	SetFinish(handler FinishHandler)
	Scheme() string
	Assert(i interface{}) error
}
