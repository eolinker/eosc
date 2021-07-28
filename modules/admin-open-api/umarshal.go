package admin_open_api

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/ghodss/yaml"
	"net/http"
)

type JsonData []byte

func (j JsonData) UnMarshal(v interface{}) error {
	return json.Unmarshal(j,&v)
}

func (j JsonData) Marshal() ([]byte, error) {
	return j,nil
}


type XMLData []byte

type YamlData []byte

func (y YamlData) UnMarshal(v interface{}) error {
	return yaml.Unmarshal(y,v)
}

func (y YamlData) Marshal() ([]byte, error) {
	return y,nil
}


type MarshalFactory interface {
	GetData(req *http.Request)eosc.IData
}

type openApiMarshal struct {

}

func (o *openApiMarshal) GetData(req *http.Request) eosc.IData {
	panic("implement me")
}
