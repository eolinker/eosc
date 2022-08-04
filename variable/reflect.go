package variable

import (
	"reflect"
)

// RecurseReflect 递归反射值给对象
func RecurseReflect(originVal reflect.Value, targetValue reflect.Value, variable map[string]string) error {
	switch originVal.Kind() {
	case reflect.Interface:
		return interfaceSet(originVal, targetValue, variable)
	case reflect.Map:
		return mapSet(originVal, targetValue, variable)
	case reflect.String:
		return stringSet(originVal, targetValue, variable)
	case reflect.Array, reflect.Slice:
		return arraySet(originVal, targetValue, variable)
	}
	return nil
}
