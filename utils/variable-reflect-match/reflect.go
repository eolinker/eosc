package main

import (
	"fmt"
	"reflect"
)

// recurseReflect 递归反射值给对象
func recurseReflect(originVal reflect.Value, targetValue reflect.Value, variable map[string]string) error {
	fmt.Println("origin kind:", originVal.Kind())
	switch originVal.Kind() {
	case reflect.Interface:
		{
			err := interfaceDeal(originVal, targetValue, variable)
			if err != nil {
				return err
			}
			if targetValue.Kind() == reflect.Ptr {
				fmt.Println("interface target value", targetValue.Elem())
			} else {
				fmt.Println("interface target value", targetValue)
			}
		}
	case reflect.Map:
		{
			err := mapDeal(originVal, targetValue, variable)
			if err != nil {
				return err
			}
			if targetValue.Kind() == reflect.Ptr {
				fmt.Println("map target value", targetValue.Elem())
			} else {
				fmt.Println("map target value", targetValue)
			}
		}
	case reflect.String:
		return stringSet(originVal.String(), targetValue, variable)
	}
	return nil
}
