package process_admin

import "encoding/json"

type configProxy struct {
	data []byte
}

func (c *configProxy) UnmarshalJSON(bytes []byte) error {
	c.data = bytes
	return nil
}

func splitConfig(data []byte) []IData {
	var ps []*configProxy

	err := json.Unmarshal(data, &ps)
	if err != nil {
		return nil
	}
	r := make([]IData, 0, len(ps))
	for _, v := range ps {
		r = append(r, JsonData(v.data))
	}
	return r
}
