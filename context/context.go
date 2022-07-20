package context

import (
	"context"
)

type CompleteHandler interface {
	Complete(ctx Context) error
}

type FinishHandler interface {
	Finish(ctx Context) error
}
type ContextR interface {
	RequestId() string
	Context() context.Context
	Value(key interface{}) interface{}
	WithValue(key, val interface{})
	Complete() CompleteHandler
	SetCompleteHandler(handler CompleteHandler)
	Finish() FinishHandler
	SetFinish(handler FinishHandler)
	Scheme() string
}

type Context interface {
	ContextR
	Assert(i interface{}) error
}
