package http_service

import "github.com/eolinker/eosc"

var (
	FilterSkillName = eosc.TypeNameOf((*IFilter)(nil))
)

type IFilter interface {
	DoFilter(ctx IHttpContext, next IChain) (err error)
}
type IChain interface {
	DoChain(ctx IHttpContext) error
}
