package http

type IFilter interface {
	DoFilter(ctx IHttpContext, endpoint IEndpoint) (err error)
}
