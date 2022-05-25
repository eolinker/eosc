package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

func professionConfig(oc map[string][]byte) []*eosc.ProfessionConfig {
	if oc == nil || len(oc) == 0 {
		return nil
	}
	configs := make([]*eosc.ProfessionConfig, 0, len(oc))
	for _, v := range oc {
		c := new(eosc.ProfessionConfig)
		err := json.Unmarshal(v, c)
		if err != nil {
			log.Error("read profession config:", err)
			continue
		}
		configs = append(configs, c)
	}

	return configs
}
