package eocontext

import "github.com/eolinker/eosc/utils/config"

var (
	FilterSkillName = config.TypeNameOf((*IFilter)(nil))
)

type IFilter interface {
	DoFilter(ctx EoContext, next IChain) (err error)
	Destroy()
}
type IChain interface {
	DoChain(ctx EoContext) error
	Destroy()
}
