package dubbo_context

import "github.com/eolinker/eosc/eocontext"

func Assert(ctx eocontext.EoContext) (IDubboContext, error) {
	var dubboContext IDubboContext
	err := ctx.Assert(&dubboContext)
	return dubboContext, err
}
