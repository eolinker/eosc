package variable

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc"
	"reflect"
	"strconv"
)

var (
	ErrorVariableNotFound = errors.New("variable not found")
	ErrorUnsupportedKind  = errors.New("unsupported kind")
)

var (
	methodName = "Reset"
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
	if targetVal.Type().Implements(reflect.TypeOf((*eosc.IPluginReset)(nil)).Elem()) {
		// 判断是否实现IPluginReset接口
		f := targetVal.MethodByName(methodName)
		if f.IsValid() {
			// 判断是否实现Reset方法，如果实现，则执行赋值操作
			if targetVal.Kind() == reflect.Ptr {
				targetVal = targetVal.Elem()
			}
			vs := f.Call([]reflect.Value{reflect.ValueOf(originVal.Elem()), reflect.ValueOf(targetVal), reflect.ValueOf(variable)})
			if len(vs) > 0 {
				err, ok := vs[0].Interface().(error)
				if ok {
					return err
				}
				return nil
			}
		}
	}
	switch originVal.Elem().Kind() {
	case reflect.Map:
		return mapSet(originVal.Elem(), targetVal, variable)
	case reflect.Array, reflect.Slice:
		return arraySet(originVal.Elem(), targetVal, variable)
	case reflect.String:
		return stringSet(originVal.Elem(), targetVal, variable)
	case reflect.Float64:
		return float64Set(originVal.Elem(), targetVal)
	case reflect.Bool:
		return boolSet(originVal.Elem(), targetVal)
	default:
		fmt.Println("interface deal", "kind", originVal.Elem().Kind())
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
		err := RecurseReflect(indexValue, newValue, variable)
		if err != nil {
			return err
		}
		newSlice = reflect.Append(newSlice, newValue.Elem())
	}
	targetVal.Set(newSlice)
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
			targetType := targetVal.Type()
			newMap := reflect.MakeMap(targetType)
			for _, key := range originVal.MapKeys() {
				newKey := reflect.New(targetType.Key())
				err := RecurseReflect(key, newKey, variable)
				if err != nil {
					return err
				}
				value := originVal.MapIndex(key)
				newValue := reflect.New(targetType.Elem())
				err = RecurseReflect(value, newValue, variable)
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
		err := RecurseReflect(value, fieldValue, variable)
		if err != nil {
			return err
		}
		targetVal.Field(i).Set(fieldValue.Elem())
	}
	return nil
}
