package line

import "github.com/eolinker/eosc/formatter"

var Name = "line"

func init() {
	formatter.Register(Name, NewFactory())
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(cfg formatter.Config) (formatter.IFormatter, error) {
	return NewLine(cfg)
}
