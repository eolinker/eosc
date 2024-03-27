package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/marshal"
)

type configProxy struct {
	data []byte
}

func (c *configProxy) UnmarshalJSON(bytes []byte) error {
	c.data = bytes
	return nil
}

func splitConfig(data []byte) []admin.IData {
	var ps []*configProxy

	err := json.Unmarshal(data, &ps)
	if err != nil {
		return nil
	}
	r := make([]admin.IData, 0, len(ps))
	for _, v := range ps {
		r = append(r, marshal.JsonData(v.data))
	}
	return r
}
