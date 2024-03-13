package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/process-admin/workers"
)

type configProxy struct {
	data []byte
}

func (c *configProxy) UnmarshalJSON(bytes []byte) error {
	c.data = bytes
	return nil
}

func splitConfig(data []byte) []workers.IData {
	var ps []*configProxy

	err := json.Unmarshal(data, &ps)
	if err != nil {
		return nil
	}
	r := make([]workers.IData, 0, len(ps))
	for _, v := range ps {
		r = append(r, marshal.JsonData(v.data))
	}
	return r
}
