package model

import "encoding/json"

type Variables map[string]string

func (v Variables) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &v)
}
