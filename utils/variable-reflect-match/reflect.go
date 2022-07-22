package main

import (
	"reflect"
)

// recurseReflect 递归反射值给对象
func recurseReflect(originVal reflect.Value, targetValue reflect.Value, variable map[string]string, name string) error {
	switch originVal.Kind() {
	case reflect.Interface:
		{
			err := interfaceDeal(originVal, targetValue, variable)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		{
			err := mapDeal(originVal, targetValue, variable, name)
			if err != nil {
				return err
			}
		}
	case reflect.String:
		return stringSet(originVal.String(), targetValue, variable)
	}
	return nil
}
