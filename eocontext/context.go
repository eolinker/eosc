package eocontext

import (
	"context"
	"errors"
	"net"
	"time"
)

var ErrEoCtxUnCloneable = errors.New("EoContext is UnCloneable. ")

type CompleteHandler interface {
	Complete(ctx EoContext) error
}

type FinishHandler interface {
	Finish(ctx EoContext) error
}

type EoContext interface {
	RequestId() string
	AcceptTime() time.Time
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
	GetBalance() BalanceHandler
	SetBalance(handler BalanceHandler)
	GetUpstreamHostHandler() UpstreamHostHandler
	SetUpstreamHostHandler(handler UpstreamHostHandler)

	RealIP() string
	LocalIP() net.IP
	LocalAddr() net.Addr
	LocalPort() int

	IsCloneable() bool
	Clone() (EoContext, error)
}
