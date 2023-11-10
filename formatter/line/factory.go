package line

import (
	"github.com/eolinker/eosc"
)

const Name = "line"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(cfg eosc.FormatterConfig, extendCfg ...interface{}) (eosc.IFormatter, error) {
	return NewLine(cfg)
}
