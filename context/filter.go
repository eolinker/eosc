package context

import "github.com/eolinker/eosc/utils/config"

var (
	FilterSkillName = config.TypeNameOf((*IFilter)(nil))
)

type IFilter interface {
	DoFilter(ctx Context, next IChain) (err error)
	Destroy()
}
type IChain interface {
	DoChain(ctx Context) error
	Destroy()
}
