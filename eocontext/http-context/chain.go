package http_context

import (
	"github.com/eolinker/eosc/eocontext"
)

func Assert(ctx eocontext.EoContext) (IHttpContext, error) {
	var httpContext IHttpContext
	err := ctx.Assert(&httpContext)
	return httpContext, err
}
