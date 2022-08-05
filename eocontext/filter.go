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

type Filters []IFilter

func (fs Filters) DoChain(ctx EoContext) error {
	if len(fs) > 0 {
		f := fs[0]
		next := fs[1:]

		return f.DoFilter(ctx, next)
	}

	return nil
}

func (fs Filters) Destroy() {
	for _, f := range fs {
		f.Destroy()
	}
}
