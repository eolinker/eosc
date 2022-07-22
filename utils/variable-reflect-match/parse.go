package main

import (
	"encoding/json"
	"errors"
	"reflect"
)

func NewParse(typ reflect.Type, variable map[string]string) (*Parse, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("error struct")
	}
	return &Parse{typ: typ, variable: variable}, nil
}

type Parse struct {
	variable map[string]string
	origin   interface{}
	typ      reflect.Type
}

func (p *Parse) String() string {
	b, _ := json.Marshal(p.origin)
	return string(b)
}

func (p *Parse) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &p.origin)
	if err != nil {
		return err
	}
	return nil
}
