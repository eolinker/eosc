package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type RequireId = eosc.RequireId

var (
	_RequireTypeName      = TypeNameOf(RequireId(""))
	_RequireSliceTypeName = TypeNameOf([]RequireId{})
)

func TypeNameOf(v interface{}) string {

	return TypeName(reflect.TypeOf(v))
}

func TypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return fmt.Sprint("*", TypeName(t.Elem()))
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.String())
}

func CheckConfig(v interface{}, workers eosc.IWorkers) (map[RequireId]eosc.IWorker, error) {

	value := reflect.ValueOf(v)
	ws, err := checkConfig(value, workers)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", TypeNameOf(v), err)
	}
	if ws == nil {
		ws = make(map[RequireId]eosc.IWorker)
	}

	return ws, nil

}

func checkConfig(v reflect.Value, workers eosc.IWorkers) (map[RequireId]eosc.IWorker, error) {
	kind := v.Kind()
	switch kind {
	case reflect.Ptr:
		if v.IsNil() {
			return nil, nil
		}
		return checkConfig(v.Elem(), workers)
	case reflect.Struct:
		t := v.Type()
		n := v.NumField()
		requires := make(map[RequireId]eosc.IWorker)
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
		requires := make(map[RequireId]eosc.IWorker)
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
		requires := make(map[RequireId]eosc.IWorker)

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
	//return nil, eosc.ErrorConfigFieldUnknown
}

func checkField(f reflect.StructField, v reflect.Value, workers eosc.IWorkers) (map[RequireId]eosc.IWorker, error) {

	typeName := TypeName(v.Type())
	switch typeName {
	case _RequireTypeName:
		{
			id, _ := url.PathUnescape(v.String())
			if id == "" {
				require, has := f.Tag.Lookup("required")
				if !has || strings.ToLower(require) != "false" {
					return nil, fmt.Errorf("%s:%w", f.Name, eosc.ErrorRequire)
				}
				return nil, nil
			}

			target, has := workers.Get(id)
			if !has || target == nil {
				require, has := f.Tag.Lookup("required")
				if !has || strings.ToLower(require) != "false" {
					return nil, fmt.Errorf("required %s:%w", id, eosc.ErrorWorkerNotExits)
				}
				return nil, nil
			}

			skill, has := f.Tag.Lookup("skill")
			if !has {
				return nil, fmt.Errorf("field %s type %s :%w", f.Name, typeName, eosc.ErrorNotGetSillForRequire)
			}
			log.DebugF("check skill:%s on %s:%v", skill, id, target)
			if !target.CheckSkill(skill) {
				return nil, fmt.Errorf(" %s type %s value %s:%w", f.Name, typeName, id, eosc.ErrorTargetNotImplementSkill)
			}
			return map[RequireId]eosc.IWorker{RequireId(id): target}, nil
		}
	case _RequireSliceTypeName:
		{
			skill, has := f.Tag.Lookup("skill")
			if !has {
				return nil, fmt.Errorf("field %s type %s :%w", f.Name, typeName, eosc.ErrorNotGetSillForRequire)
			}
			require, requireHas := f.Tag.Lookup("require")

			n := v.Len()
			requires := make(map[RequireId]eosc.IWorker)
			for i := 0; i < n; i++ {
				id := v.Index(i).String()
				if id == "" {
					continue
				}
				target, has := workers.Get(id)
				if !has {
					if !requireHas || strings.ToLower(require) != "false" {
						return nil, fmt.Errorf("require %s:%w", id, eosc.ErrorWorkerNotExits)
					}
				}
				if !target.CheckSkill(skill) {
					return nil, fmt.Errorf(" %s type %s value %s:%w", f.Name, typeName, id, eosc.ErrorTargetNotImplementSkill)
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

func merge(dist map[RequireId]eosc.IWorker, source map[RequireId]eosc.IWorker) map[RequireId]eosc.IWorker {
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
