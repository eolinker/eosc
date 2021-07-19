package admin_open_api

import "encoding/json"

type JsonData []byte

func (j JsonData) UnMarshal(v interface{}) error {
	return json.Unmarshal(j,&v)
}

func (j JsonData) Marshal() ([]byte, error) {
	panic("implement me")
}


type XMLData []byte

