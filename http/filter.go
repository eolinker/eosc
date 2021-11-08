package http

type IFilter interface {
	DoFilter(ctx IHttpContext, endpoint IEndpoint, next IFilterChain) (err error)
}
type IFilterChain interface {
	DoFilter(ctx IHttpContext, endpoint IEndpoint) error
}
type IChain interface {
	IFilterChain
	Append(filter IFilter)
	Insert(filter IFilter)
}
