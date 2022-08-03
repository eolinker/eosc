package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrorVariableNotFound = errors.New("variable not found")
	ErrorUnsupportedKind  = errors.New("unsupported kind")
)

func stringSet(value reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	builder := NewBuilder(value.String())
	val, success := builder.Replace(variable)
	if !success {
		return ErrorVariableNotFound
	}
	switch targetVal.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("string set parse int error: %w", err)
		}
		targetVal.SetInt(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("string set parse bool error: %w", err)
		}
		targetVal.SetBool(v)
	case reflect.String:
		targetVal.SetString(val)
	default:
		return fmt.Errorf("%w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	return nil
}

func interfaceSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	value := originVal.Elem()
	switch value.Kind() {
	case reflect.Map:
		return mapSet(value, targetVal, variable)
	case reflect.Array, reflect.Slice:
		return arraySet(value, targetVal, variable)
	case reflect.String:
		return stringSet(value, targetVal, variable)
	case reflect.Float64:
		return float64Set(value, targetVal)
	case reflect.Bool:
		return boolSet(value, targetVal)
	default:
		fmt.Println("interface deal", "kind", value.Kind())
	}
	return nil
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
	case reflect.Int:
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

func mapSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	if originVal.Kind() != reflect.Map {
		return fmt.Errorf("map deal %w %s", ErrorUnsupportedKind, originVal.Kind())
	}
	if targetVal.Kind() == reflect.Ptr {
		if !targetVal.Elem().IsValid() {
			targetType := targetVal.Type()
			newVal := reflect.New(targetType.Elem())
			targetVal.Set(newVal)
		}
		targetVal = targetVal.Elem()
	}
	switch targetVal.Kind() {
	case reflect.Struct:
		{
			return structSet(originVal, targetVal, variable)
		}
	case reflect.Map:
		{
			if targetVal.Type() == reflect.TypeOf(PluginMap{}) {
				return PluginMapSet(originVal, targetVal, variable)
			}
			targetType := targetVal.Type()
			newMap := reflect.MakeMap(targetType)
			for _, key := range originVal.MapKeys() {
				newKey := reflect.New(targetType.Key())
				err := recurseReflect(key, newKey, variable)
				if err != nil {
					return err
				}
				value := originVal.MapIndex(key)
				newValue := reflect.New(targetType.Elem())
				err = recurseReflect(value, newValue, variable)
				if err != nil {
					return err
				}
				newMap.SetMapIndex(newKey.Elem(), newValue.Elem())
			}
			targetVal.Set(newMap)
		}
	case reflect.Ptr:
		{
			err := mapSet(originVal, targetVal, variable)
			if err != nil {
				return err
			}
		}

	default:
		{
			fmt.Println("type", targetVal.Type(), "kind", targetVal.Kind())
		}
	}
	return nil
}

func structSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	targetType := targetVal.Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := reflect.New(field.Type)
		tag := field.Tag.Get("json")
		value := originVal.MapIndex(reflect.ValueOf(tag))
		err := recurseReflect(value, fieldValue, variable)
		if err != nil {
			return err
		}
		targetVal.Field(i).Set(fieldValue.Elem())
	}
	return nil
}
