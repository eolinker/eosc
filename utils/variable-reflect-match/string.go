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

func stringSet(value string, targetVal reflect.Value, variable map[string]string) error {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	builder := NewBuilder(value)
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
