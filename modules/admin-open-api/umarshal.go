package admin_open_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
)

var (
	ErrorUnknownContentType = errors.New("unknown content type")
)

type IData interface {
	UnMarshal(v interface{}) error
	Encode() ([]byte, error)
}
type JsonData []byte

func (j JsonData) Encode() ([]byte, error) {
	return j, nil
}

func (j JsonData) UnMarshal(v interface{}) error {
	return json.Unmarshal(j, &v)
}

func (j JsonData) Marshal() ([]byte, error) {
	return j, nil
}

type XMLData []byte

type YamlData []byte

func (y YamlData) Encode() ([]byte, error) {
	v := make(map[string]interface{})
	err := y.UnMarshal(&v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (y YamlData) UnMarshal(v interface{}) error {
	return yaml.Unmarshal(y, v)
}

func (y YamlData) Marshal() ([]byte, error) {
	return y, nil
}

func GetData(req *http.Request) (IData, error) {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("content-type"))
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(mediaType) {
	case "application/json":
		data, e := ioutil.ReadAll(req.Body)
		if e != nil {
			return nil, e
		}
		req.Body.Close()

		return JsonData(data), nil
	case "application/yaml":
		data, e := ioutil.ReadAll(req.Body)
		if e != nil {
			return nil, e
		}
		req.Body.Close()

		return YamlData(data), nil

	}

	return nil, fmt.Errorf("%s:%w", mediaType, ErrorUnknownContentType)
}
