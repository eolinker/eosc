package main

import (
	"fmt"
	"reflect"
)

func arraySet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	if originVal.Kind() != reflect.Slice && originVal.Kind() != reflect.Array {
		return fmt.Errorf("origin error: %w %s", ErrorUnsupportedKind, originVal.Kind())
	}
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	if targetVal.Kind() != reflect.Slice {
		return fmt.Errorf("target error %w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	newSlice := reflect.MakeSlice(targetVal.Type(), 0, originVal.Cap())
	for j := 0; j < originVal.Len(); j++ {
		indexValue := originVal.Index(j)
		newValue := reflect.New(targetVal.Type().Elem())
		err := recurseReflect(indexValue, newValue, variable, "")
		if err != nil {
			return err
		}
		newSlice = reflect.Append(newSlice, newValue.Elem())
	}
	targetVal.Set(newSlice)
	return nil
}
