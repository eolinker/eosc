package http_context

import "github.com/eolinker/eosc/context"

func Assert(ctx context.Context) (IHttpContext, error) {
	var httpContext IHttpContext
	err := ctx.Assert(&httpContext)
	return httpContext, err
}
