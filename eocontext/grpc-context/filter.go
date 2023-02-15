package grpc_context

import (
	"github.com/eolinker/eosc/eocontext"
	"github.com/eolinker/eosc/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*GrpcFilter)(nil))
)

type GrpcFilter interface {
	DoGrpcFilter(ctx IGrpcContext, next eocontext.IChain) (err error)
}

func DoGrpcFilter(filter GrpcFilter, ctx eocontext.EoContext, next eocontext.IChain) (err error) {
	grpcContext, err := Assert(ctx)
	if err == nil {
		return filter.DoGrpcFilter(grpcContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}
