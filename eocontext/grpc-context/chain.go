package grpc_context

import (
	"github.com/eolinker/eosc/eocontext"
)

func Assert(ctx eocontext.EoContext) (IGrpcContext, error) {
	var grpcContext IGrpcContext
	err := ctx.Assert(&grpcContext)
	return grpcContext, err
}
