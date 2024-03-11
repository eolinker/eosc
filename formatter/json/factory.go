package json

import (
	"encoding/json"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
)

const Name = "json"

type Factory struct {
}

func (f *Factory) Create(cfg eosc.FormatterConfig, extendCfg ...interface{}) (eosc.IFormatter, error) {
	var ctRs []contentResize
	if len(extendCfg) > 0 {
		err := json.Unmarshal(extendCfg[0].([]byte), &ctRs)
		if err != nil {
			log.Errorf("json formatter extend config error: %s", err)
		}
	}
	return NewFormatter(cfg, ctRs)
}

func NewFactory() *Factory {
	return &Factory{}
}
