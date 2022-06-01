package http_service

import (
	"github.com/eolinker/eosc/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*IFilter)(nil))
)

type IFilter interface {
	DoFilter(ctx IHttpContext, next IChain) (err error)
	Destroy()
}
type IChain interface {
	DoChain(ctx IHttpContext) error
	Destroy()
}
