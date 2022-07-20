package http_context

import (
	"github.com/eolinker/eosc/context"
	"github.com/eolinker/eosc/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*HttpFilter)(nil))
)

type HttpFilter interface {
	DoHttpFilter(ctx IHttpContext, next context.IChain) (err error)
}

func DoHttpFilter(httpFilter HttpFilter, ctx context.Context, next context.IChain) (err error) {
	httpContext, err := Assert(ctx)
	if err == nil {
		return httpFilter.DoHttpFilter(httpContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}
