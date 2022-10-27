package http_context

import (
	"github.com/eolinker/eosc/eocontext"
	"github.com/eolinker/eosc/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*HttpFilter)(nil))
)

type HttpFilter interface {
	DoHttpFilter(ctx IHttpContext, next eocontext.IChain) (err error)
}

func DoHttpFilter(httpFilter HttpFilter, ctx eocontext.EoContext, next eocontext.IChain) (err error) {
	httpContext, err := Assert(ctx)
	if err == nil {
		return httpFilter.DoHttpFilter(httpContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}

type WebsocketFilter interface {
	DoWebsocketFilter(ctx IWebsocketContext, next eocontext.IChain) error
}
