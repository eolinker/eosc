package setting

import "encoding/json"

type configProxy struct {
	data []byte
}

func (c *configProxy) UnmarshalJSON(bytes []byte) error {
	c.data = bytes
	return nil
}

func splitConfig(data []byte) [][]byte {
	var ps []*configProxy

	err := json.Unmarshal(data, &ps)
	if err != nil {
		return [][]byte{data}
	}
	r := make([][]byte, 0, len(ps))
	for _, v := range ps {
		r = append(r, v.data)
	}
	return r
}
