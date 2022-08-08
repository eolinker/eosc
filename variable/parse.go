package variable

import (
	"encoding/json"
	"reflect"
)

func NewParse(variables map[string]string, configTypes map[string]reflect.Type) (*Parse, error) {
	return &Parse{variables: variables, configTypes: configTypes}, nil
}

type Parse struct {
	variables   map[string]string
	configTypes map[string]reflect.Type
}

func (p *Parse) Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error) {
	o := newOrg(typ, p.variables, p.configTypes)
	err := json.Unmarshal(buf, o)
	return o.target, o.usedVariable, err
}

type org struct {
	typ          reflect.Type
	variable     map[string]string
	configTypes  map[string]reflect.Type
	target       interface{}
	usedVariable []string
}

func newOrg(typ reflect.Type, variable map[string]string, configTypes map[string]reflect.Type) *org {
	return &org{typ: typ, variable: variable, configTypes: configTypes}
}

func (o *org) UnmarshalJSON(bytes []byte) error {
	var origin interface{}
	err := json.Unmarshal(bytes, &origin)
	if err != nil {
		return err
	}
	target := reflect.New(o.typ)
	variables, err := recurseReflect(reflect.ValueOf(origin), target, o.variable, o.configTypes)
	if err != nil {
		return err
	}
	o.target = target.Elem().Interface()
	variableMap := make(map[string]bool)
	for _, v := range variables {
		if _, ok := variableMap[v]; !ok {
			variableMap[v] = true
		}
	}
	usedVariables := make([]string, 0, len(variableMap))
	for key := range variableMap {
		usedVariables = append(usedVariables, key)
	}
	o.usedVariable = usedVariables
	return nil
}
