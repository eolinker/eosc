package eosc

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type JsonData []byte

func (j JsonData) UnMarshal(v interface{}) error {
	return json.Unmarshal(j, &v)
}

func (j JsonData) Marshal() ([]byte, error) {
	return j, nil
}

type XMLData []byte

type YamlData []byte

func (y YamlData) UnMarshal(v interface{}) error {
	return yaml.Unmarshal(y, v)
}

func (y YamlData) Marshal() ([]byte, error) {
	return y, nil
}
