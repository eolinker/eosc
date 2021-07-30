package eosc

import (
	"fmt"
	"reflect"
	"strings"
)

type RequireId string

var (
	_RequireTypeName      = TypeNameOf(RequireId(""))
	_RequireSliceTypeName = TypeNameOf([]RequireId{})
)

func TypeNameOf(v interface{}) string {

	return TypeName(reflect.TypeOf(v))
}

//func TypeNameOfValue(v reflect.Value) string {
//	 if v.Kind() == reflect.Ptr{
//	 	return TypeNameOfValue(v.Elem())
//	 }
//	 return TypeName(v.Type())
//}
func TypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return TypeName(t.Elem())
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.String())
}

func CheckConfig(v interface{}, workers IWorkers) (map[RequireId]interface{}, error) {

	value := reflect.ValueOf(v)
	ws, err := checkConfig(value, workers)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", TypeNameOf(v), err)
	}
	if ws == nil {
		ws = make(map[RequireId]interface{})
	}

	return ws, nil

}

func checkConfig(v reflect.Value, workers IWorkers) (map[RequireId]interface{}, error) {
	kind := v.Kind()
	switch kind {
	case reflect.Ptr:
		if v.IsNil() {
			return nil, ErrorConfigIsNil
		}
		return checkConfig(v.Elem(), workers)
	case reflect.Struct:
		t := v.Type()
		n := v.NumField()
		requires := make(map[RequireId]interface{})
		for i := 0; i < n; i++ {
			if ws, err := checkField(t.Field(i), v.Field(i), workers); err != nil {
				return nil, err
			} else {
				requires = merge(requires, ws)
			}
		}
		return requires, nil
	case reflect.Slice:
		n := v.Len()
		requires := make(map[RequireId]interface{})
		for i := 0; i < n; i++ {
			if ws, err := checkConfig(v.Index(i), workers); err != nil {
				return nil, err
			} else {
				requires = merge(requires, ws)
			}
		}
		return requires, nil
	case reflect.Map:
		it := v.MapRange()
		requires := make(map[RequireId]interface{})

		for it.Next() {
			if ws, err := checkConfig(it.Value(), workers); err != nil {
				return nil, err
			} else {
				requires = merge(requires, ws)
			}
		}
		return requires, nil
	default:
		return nil, nil
	}
	return nil, ErrorConfigFieldUnknown
}

func checkField(f reflect.StructField, v reflect.Value, workers IWorkers) (map[RequireId]interface{}, error) {

	typeName := TypeName(v.Type())
	switch typeName {
	case _RequireTypeName:
		{
			id := v.String()
			if id == "" {
				return nil, fmt.Errorf("%s:%w", f.Name, ErrorRequire)
			}

			target, has := workers.Get(id)
			if !has {
				require, has := f.Tag.Lookup("require")
				if !has || strings.ToLower(require) != "false" {
					return nil, fmt.Errorf("require %s:%w", id, ErrorWorkerNotExits)
				}
				return nil, nil
			}

			skill, has := f.Tag.Lookup("skill")
			if !has {
				return nil, fmt.Errorf("field %s type %s :%w", f.Name, typeName, ErrorNotGetSillForRequire)
			}
			if !target.CheckSkill(skill) {
				return nil, fmt.Errorf(" %s type %s value %s:%w", f.Name, typeName, id, ErrorTargetNotImplementSkill)
			}
			return map[RequireId]interface{}{RequireId(id): target}, nil
		}
	case _RequireSliceTypeName:
		{
			skill, has := f.Tag.Lookup("skill")
			if !has {
				return nil, fmt.Errorf("field %s type %s :%w", f.Name, typeName, ErrorNotGetSillForRequire)
			}
			require, requireHas := f.Tag.Lookup("require")

			n := v.Len()
			requires := make(map[RequireId]interface{})
			for i := 0; i < n; i++ {
				id := v.Index(i).String()
				if id == "" {
					continue
				}
				target, has := workers.Get(id)
				if !has {
					if !requireHas || strings.ToLower(require) != "false" {
						return nil, fmt.Errorf("require %s:%w", id, ErrorWorkerNotExits)
					}
				}
				if !target.CheckSkill(skill) {
					return nil, fmt.Errorf(" %s type %s value %s:%w", f.Name, typeName, id, ErrorTargetNotImplementSkill)
				}
				requires[RequireId(id)] = target
			}
			return requires, nil
		}
	default:
		{
			return checkConfig(v, workers)
		}
	}
}

func merge(dist map[RequireId]interface{}, source map[RequireId]interface{}) map[RequireId]interface{} {
	if dist == nil && source == nil {
		return nil
	}
	if source == nil {
		return dist
	}
	if dist == nil {
		return source
	}
	for k, v := range source {
		dist[k] = v
	}
	return dist
}

func newConfig(t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
