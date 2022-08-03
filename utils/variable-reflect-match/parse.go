package main

import (
	"encoding/json"
	"reflect"
)

func NewParse(variable map[string]string) (*Parse, error) {
	return &Parse{variable: variable}, nil
}

type Parse struct {
	variable map[string]string
}

func (p *Parse) Unmarshal(buf []byte, typ reflect.Type) (interface{}, error) {
	o := newOrg(typ, p.variable)
	err := json.Unmarshal(buf, o)
	return o.target, err
}

type org struct {
	typ      reflect.Type
	variable map[string]string
	target   interface{}
}

func newOrg(typ reflect.Type, variable map[string]string) *org {
	return &org{typ: typ, variable: variable}
}

func (o *org) UnmarshalJSON(bytes []byte) error {
	var origin interface{}
	err := json.Unmarshal(bytes, &origin)
	if err != nil {
		return err
	}
	target := reflect.New(o.typ)
	err = recurseReflect(reflect.ValueOf(origin), target, o.variable)
	if err != nil {
		return err
	}
	o.target = target.Elem().Interface()
	return nil
}
