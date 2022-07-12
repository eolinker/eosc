package context

type NodeAddress = string
type APP interface {
	GetNode(id string) (NodeAddress, bool)
	All() map[string]NodeAddress
}
type LoadBalance interface {
	Select(app APP) (NodeAddress, error)
}
type DoHandler interface {
	DO() error
}
type FinishHandler interface {
	Finish() error
}
type Context interface {
	LoadBalance() LoadBalance
	SetLoadBalance(balance LoadBalance)
	DO() DoHandler
	SetDoHandler(handler DoHandler)
	Finish() FinishHandler
	SetFinish(handler FinishHandler)
	Assert(i interface{}) error
	Scheme() string
}
