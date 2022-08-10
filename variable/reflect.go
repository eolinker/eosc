package variable

import (
	"github.com/eolinker/eosc"
	"reflect"
)

var (
	methodName = "Reset"
)

func RecurseReflect(originVal reflect.Value, targetVal reflect.Value, variables map[string]string) ([]string, error) {
	return recurseReflect(originVal, targetVal, variables)
}

// recurseReflect 递归反射值给对象
func recurseReflect(originVal reflect.Value, targetVal reflect.Value, variables map[string]string) ([]string, error) {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	usedVariables := make([]string, 0, len(variables))
	var used []string
	var err error
	switch originVal.Kind() {
	case reflect.Interface:
		if targetVal.Type().Implements(reflect.TypeOf((*eosc.IPluginReset)(nil)).Elem()) {
			// 判断是否实现Reset方法，如果实现，则执行赋值操作
			f := targetVal.MethodByName(methodName)
			if f.IsValid() {
				if targetVal.Kind() == reflect.Ptr {
					// 空指针赋值
					if !targetVal.Elem().IsValid() {
						// 当elem非法类型时，进行对象赋值
						newVal := reflect.New(targetVal.Type().Elem())
						targetVal.Set(newVal)
					}
					targetVal = targetVal.Elem()
				}

				vs := f.Call([]reflect.Value{reflect.ValueOf(originVal.Elem()), reflect.ValueOf(targetVal), reflect.ValueOf(variables)})
				err, ok := vs[1].Interface().(error)
				if ok {
					return nil, err
				}
				used, ok = vs[0].Interface().([]string)
				if !ok {
					return nil, nil
				}
				usedVariables = append(usedVariables, used...)
				return usedVariables, nil
			}
		}
		used, err = interfaceSet(originVal, targetVal, variables)
	case reflect.Map:
		used, err = mapSet(originVal, targetVal, variables)
	case reflect.String:
		used, err = stringSet(originVal, targetVal, variables)
	case reflect.Array, reflect.Slice:
		used, err = arraySet(originVal, targetVal, variables)
	}
	usedVariables = append(usedVariables, used...)
	return usedVariables, err
}
