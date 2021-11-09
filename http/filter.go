package http

type IFilter interface {
	DoFilter(ctx IHttpContext, next IChain) (err error)
}
type IChain interface {
	DoFilter(ctx IHttpContext) error
}
