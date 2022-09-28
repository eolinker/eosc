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
type IChainPro interface {
	Destroy()
	Chain(ctx EoContext, append ...IFilter) error
}
type Filters []IFilter

func (fs Filters) DoChain(ctx EoContext) error {
	if len(fs) > 0 {
		f := fs[0]
		next := fs[1:]
		return f.DoFilter(ctx, next)
	}

	//ctx.GetComplete().Complete(ctx)
	return nil
}

func (fs Filters) Destroy() {
	for _, f := range fs {
		f.Destroy()
	}
}

func DoChain(ctx EoContext, orgfilter Filters, append ...IFilter) error {
	fs := orgfilter
	fl := len(fs)
	al := len(append)
	if fl == 0 && al == 0 {
		return nil
	}
	if fl == 0 {
		return Filters(append).DoChain(ctx)
	}
	if al == 0 {
		return fs.DoChain(ctx)
	}

	tp := make(Filters, fl+al)
	copy(tp, fs)
	copy(tp[fl:], append)
	return tp.DoChain(ctx)
}

type _FilterChain struct {
	chain IChain
}

func (c *_FilterChain) DoFilter(ctx EoContext, next IChain) (err error) {
	return c.chain.DoChain(ctx)
}

func (c *_FilterChain) Destroy() {
	c.chain.Destroy()
}

func ToFilter(chain IChain) IFilter {
	return &_FilterChain{chain: chain}
}
