package variable

import (
	"encoding/json"
	"reflect"
)

func NewParse(variables IVariable) *Parse {
	return &Parse{variables: variables}
}

type Parse struct {
	variables IVariable
}

func (p *Parse) Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error) {
	o := newOrg(typ, p.variables)
	err := json.Unmarshal(buf, o)
	return o.target, o.usedVariable, err
}

type org struct {
	typ          reflect.Type
	variable     IVariable
	target       interface{}
	usedVariable []string
}

func newOrg(typ reflect.Type, variable IVariable) *org {
	return &org{typ: typ, variable: variable}
}

func (o *org) UnmarshalJSON(bytes []byte) error {
	var origin interface{}
	err := json.Unmarshal(bytes, &origin)
	if err != nil {
		return err
	}
	target := reflect.New(o.typ)
	variables, err := recurseReflect(reflect.ValueOf(origin), target, o.variable)
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
