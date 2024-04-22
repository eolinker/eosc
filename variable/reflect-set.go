package variable

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var (
	ErrorVariableNotFound = errors.New("data not found")
	ErrorUnsupportedKind  = errors.New("unsupported kind")
)

func stringSet(value reflect.Value, targetVal reflect.Value, variables eosc.IVariable) ([]string, error) {
	if targetVal.Kind() == reflect.Ptr {
		return stringSet(value, targetVal.Elem(), variables)
	}
	builder := NewBuilder(value.String())
	val, useVariables, success := builder.Replace(variables)
	if !success {

		return nil, fmt.Errorf("<%s> %w", value.String(), ErrorVariableNotFound)
	}
	switch targetVal.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("<%s> string set parse int error: %w", value.String(), err)
		}
		targetVal.SetInt(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("<%s> string set parse bool error: %w", value.String(), err)
		}
		targetVal.SetBool(v)
	case reflect.String:
		targetVal.SetString(val)
	default:
		return nil, fmt.Errorf("<%s> %w %s", value.String(), ErrorUnsupportedKind, targetVal.Kind())
	}
	return useVariables, nil
}

func interfaceSet(originVal reflect.Value, targetVal reflect.Value, variables eosc.IVariable) ([]string, error) {
	usedVariables := make([]string, 0, variables.Len())
	var used []string
	var err error
	switch originVal.Elem().Kind() {
	case reflect.Map:
		used, err = mapSet(originVal.Elem(), targetVal, variables)
	case reflect.Array, reflect.Slice:
		used, err = arraySet(originVal.Elem(), targetVal, variables)
	case reflect.String:
		used, err = stringSet(originVal.Elem(), targetVal, variables)
	case reflect.Float64:
		err = float64Set(originVal.Elem(), targetVal)
	case reflect.Bool:
		err = boolSet(originVal.Elem(), targetVal)
	default:
		err = fmt.Errorf("interface deal kind: %s", originVal.Elem().Kind().String())
		log.Error(err)
		return nil, err
	}
	usedVariables = append(usedVariables, used...)
	return usedVariables, err
}

func boolSet(originVal reflect.Value, targetVal reflect.Value) error {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	switch targetVal.Kind() {
	case reflect.String:
		targetVal.SetString(strconv.FormatBool(originVal.Bool()))
	case reflect.Bool:
		targetVal.Set(originVal)
	default:
		return fmt.Errorf("bool set error:%w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	return nil
}

func float64Set(originVal reflect.Value, targetVal reflect.Value) error {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	switch targetVal.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(fmt.Sprintf("%1.0f", originVal.Float()), 10, 64)
		if err != nil {
			return err
		}
		targetVal.SetInt(value)
	case reflect.Float64:
		targetVal.SetFloat(originVal.Float())
	case reflect.String:
		value := fmt.Sprintf("%f", originVal.Float())
		targetVal.SetString(value)
	default:
		return fmt.Errorf("float64 set error:%w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	return nil
}

func arraySet(originVal reflect.Value, targetVal reflect.Value, variables eosc.IVariable) ([]string, error) {
	if originVal.Kind() != reflect.Slice && originVal.Kind() != reflect.Array {
		return nil, fmt.Errorf("origin error: %w %s", ErrorUnsupportedKind, originVal.Kind())
	}
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	if targetVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("target error %w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	usedVariables := make([]string, 0, variables.Len())
	newSlice := reflect.MakeSlice(targetVal.Type(), 0, originVal.Cap())
	for j := 0; j < originVal.Len(); j++ {
		indexValue := originVal.Index(j)
		newValue := reflect.New(targetVal.Type().Elem())
		used, err := recurseReflect(indexValue, newValue, variables)
		if err != nil {
			return nil, fmt.Errorf("[%d]%v", j, err)
		}
		usedVariables = append(usedVariables, used...)
		newSlice = reflect.Append(newSlice, newValue.Elem())
	}
	targetVal.Set(newSlice)
	return usedVariables, nil
}

func mapSet(originVal reflect.Value, targetVal reflect.Value, variables eosc.IVariable) ([]string, error) {
	if originVal.Kind() != reflect.Map {
		return nil, fmt.Errorf(" map deal %w %s", ErrorUnsupportedKind, originVal.Kind())
	}
	if targetVal.Kind() == reflect.Ptr {
		if !targetVal.Elem().IsValid() {
			targetType := targetVal.Type()
			newVal := reflect.New(targetType.Elem())
			targetVal.Set(newVal)
		}
		targetVal = targetVal.Elem()
	}
	usedVariables := make([]string, 0, variables.Len())
	switch targetVal.Kind() {
	case reflect.Struct:
		{
			return structSet(originVal, targetVal, variables)
		}
	case reflect.Map:
		{
			targetType := targetVal.Type()
			newMap := reflect.MakeMap(targetType)
			for _, key := range originVal.MapKeys() {
				newKey := reflect.New(targetType.Key())
				_, err := recurseReflect(key, newKey, variables)
				if err != nil {
					return nil, fmt.Errorf("[%s]%v", key.String(), err)
				}
				value := originVal.MapIndex(key)
				newValue := reflect.New(targetType.Elem())
				used, err := recurseReflect(value, newValue, variables)
				if err != nil {
					return nil, fmt.Errorf("[%s]%v", key.String(), err)
				}
				usedVariables = append(usedVariables, used...)
				newMap.SetMapIndex(newKey.Elem(), newValue.Elem())
			}
			targetVal.Set(newMap)
		}
	case reflect.Ptr:
		{
			used, err := mapSet(originVal, targetVal, variables)
			if err != nil {
				return nil, err
			}
			usedVariables = append(usedVariables, used...)
		}
	default:
		{
			log.Error("type ", targetVal.Type(), " kind ", targetVal.Kind(), " ", originVal, " ", targetVal.Type().Name())
		}
	}
	return usedVariables, nil
}

func structSet(originVal reflect.Value, targetVal reflect.Value, variables eosc.IVariable) ([]string, error) {
	usedVariables := make([]string, 0, variables.Len())
	targetType := targetVal.Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)

		fieldValue := reflect.New(field.Type)
		tag := field.Tag.Get("json")
		ts := strings.Split(tag, ",")
		value := originVal.MapIndex(reflect.ValueOf(ts[0]))
		used, err := recurseReflect(value, fieldValue, variables)
		if err != nil {
			return nil, fmt.Errorf(".%s%v", field.Name, err)
		}
		usedVariables = append(usedVariables, used...)
		if targetVal.Field(i).CanSet() {
			targetVal.Field(i).Set(fieldValue.Elem())
		}
	}
	return usedVariables, nil
}
