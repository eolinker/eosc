package admin_open_api

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/ghodss/yaml"
	"mime"
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



func  GetData(req *http.Request) (eosc.IData,error) {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("content-type"))
	if err!= nil{

	}
}
