package json

import (
	"github.com/eolinker/eosc"
)

const Name = "json"

type Factory struct {
}

func (f *Factory) Create(cfg eosc.FormatterConfig) (eosc.IFormatter, error) {

	return NewFormatter(cfg)
}

func NewFactory() *Factory {
	return &Factory{}
}
