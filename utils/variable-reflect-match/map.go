package main

import (
	"fmt"
	"reflect"
)

func mapDeal(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	if originVal.Kind() != reflect.Map {
		return fmt.Errorf("map deal %w %s", ErrorUnsupportedKind, originVal.Kind())
	}
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	fmt.Println("kind is ", targetVal.Kind())
	switch targetVal.Kind() {
	case reflect.Struct:
		{
			targetType := targetVal.Type()
			for i := 0; i < targetType.NumField(); i++ {
				field := targetType.Field(i)
				fieldValue := reflect.New(field.Type)
				tag := field.Tag.Get("json")
				err := recurseReflect(originVal.MapIndex(reflect.ValueOf(tag)), fieldValue, variable)
				if err != nil {
					return err
				}
				targetVal.Field(i).Set(fieldValue.Elem())
			}
		}
	case reflect.Map:
		{
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
			fmt.Println("map deal", originVal, "kind", targetVal.Kind(), targetVal.Type())
			err := mapDeal(originVal, targetVal, variable)
			if err != nil {
				return err
			}
			fmt.Println("map deal", originVal, "value", targetVal)
		}
	case reflect.Interface:
		{
			fmt.Println("map deal interface", originVal, targetVal.Type())
		}
	}
	return nil
}
