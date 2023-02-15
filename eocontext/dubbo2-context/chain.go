package dubbo2_context

import "github.com/eolinker/eosc/eocontext"

func Assert(ctx eocontext.EoContext) (IDubbo2Context, error) {
	var dubboContext IDubbo2Context
	err := ctx.Assert(&dubboContext)
	return dubboContext, err
}
