package main

import (
	"fmt"
	"reflect"
)

func mapSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string, name string) error {
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
			return structDeal(originVal, targetVal, variable)
		}
	case reflect.Map:
		{
			targetType := targetVal.Type()
			newMap := reflect.MakeMap(targetType)
			for _, key := range originVal.MapKeys() {
				newKey := reflect.New(targetType.Key())
				err := recurseReflect(key, newKey, variable, "")
				if err != nil {
					return err
				}
				value := originVal.MapIndex(key)
				newValue := reflect.New(targetType.Elem())
				err = recurseReflect(value, newValue, variable, key.String())
				if err != nil {
					return err
				}
				newMap.SetMapIndex(newKey.Elem(), newValue.Elem())
			}
			targetVal.Set(newMap)
		}
	case reflect.Ptr:
		{
			err := mapSet(originVal, targetVal, variable, name)
			if err != nil {
				return err
			}
		}
	case reflect.Interface:
		{
			val := reflect.ValueOf(&Config{})
			newVal := reflect.New(val.Type())
			err := mapSet(originVal, newVal, variable, "")
			if err != nil {
				return err
			}
			targetVal.Set(newVal.Elem())
		}
	}
	return nil
}

func structDeal(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	targetType := targetVal.Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := reflect.New(field.Type)
		tag := field.Tag.Get("json")
		value := originVal.MapIndex(reflect.ValueOf(tag))
		err := recurseReflect(value, fieldValue, variable, "")
		if err != nil {
			return err
		}
		targetVal.Field(i).Set(fieldValue.Elem())
	}
	return nil
}
