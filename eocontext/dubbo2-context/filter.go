package dubbo2_context

import (
	"github.com/eolinker/eosc/eocontext"
	"github.com/eolinker/eosc/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*DubboFilter)(nil))
)

type DubboFilter interface {
	DoDubboFilter(ctx IDubbo2Context, next eocontext.IChain) (err error)
}

func DoDubboFilter(httpFilter DubboFilter, ctx eocontext.EoContext, next eocontext.IChain) (err error) {
	httpContext, err := Assert(ctx)
	if err == nil {
		return httpFilter.DoDubboFilter(httpContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}
