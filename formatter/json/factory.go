package json

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/formatter"
)

const Name = "json"

func init() {
	formatter.Register(Name, NewFactory())
}

type Factory struct {
}

func (f *Factory) Create(cfg eosc.FormatterConfig) (eosc.IFormatter, error) {
	return NewFormatter(cfg)
}

func NewFactory() *Factory {
	return &Factory{}
}
