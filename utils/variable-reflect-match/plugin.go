package main

import (
	"fmt"
	"reflect"
)

//PluginConfig 普通插件配置，在router、service、upstream的插件格式
type PluginConfig struct {
	Disable bool        `json:"disable"`
	Config  interface{} `json:"config"`
}

type Config struct {
	Scheme   int               `json:"scheme" label:"协议"`
	URI      string            `json:"uri" label:"URI"`
	RegexURI []string          `json:"regex_uri" label:"正则替换URI（regex_uri）"`
	Host     string            `json:"host" label:"Host"`
	Headers  map[string]string `json:"headers" label:"请求头部"`
}

type RewriteConfig struct {
	StatusCode int               `json:"status_code" label:"响应状态码" minimum:"100" description:"最小值：100"`
	Body       string            `json:"body" label:"响应内容"`
	BodyBase64 bool              `json:"body_base64" label:"是否base64加密"`
	Headers    map[string]string `json:"headers" label:"响应头部"`
	Match      *MatchConf        `json:"match" label:"匹配状态码列表"`
}

type MatchConf struct {
	Code []int `json:"code" label:"状态码" minimum:"100" description:"最小值：100"`
}

type PluginMap map[string]*PluginConfig

type IPluginConfig interface {
	Reset(originVal reflect.Value, variable map[string]string)
}

func PluginMapSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	// originVal type: map[string]interface{}
	// targetVal type: PluginMap
	targetType := targetVal.Type()
	newMap := reflect.MakeMap(targetType)
	for _, key := range originVal.MapKeys() {
		// 判断是否存在对应的插件配置
		cfgType, ok := typeMap[key.String()]
		if !ok {
			return fmt.Errorf("plugin %s not found", key.String())
		}
		value := originVal.MapIndex(key)
		newValue := reflect.New(targetType.Elem())

		err := PluginConfigSet(value, newValue, variable, cfgType)
		if err != nil {
			return err
		}
		newMap.SetMapIndex(key, newValue.Elem())
	}
	targetVal.Set(newMap)
	return nil
}

func PluginConfigSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string, cfgType reflect.Type) error {
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
			targetType := targetVal.Type()
			for i := 0; i < targetType.NumField(); i++ {
				field := targetType.Field(i)
				var fieldValue reflect.Value
				switch field.Type.Kind() {
				case reflect.Interface:
					if cfgType.Kind() == reflect.Ptr {
						cfgType = cfgType.Elem()
					}
					fieldValue = reflect.New(cfgType)
				default:
					fieldValue = reflect.New(field.Type)
				}

				var value reflect.Value
				switch originVal.Elem().Kind() {
				case reflect.Map:
					{
						tag := field.Tag.Get("json")
						value = originVal.Elem().MapIndex(reflect.ValueOf(tag))
					}
				default:
					value = originVal.Elem()
				}

				err := recurseReflect(value, fieldValue, variable)
				if err != nil {
					return err
				}
				targetVal.Field(i).Set(fieldValue.Elem())
			}
		}
	case reflect.Ptr:
		return PluginConfigSet(originVal, targetVal, variable, cfgType)
	}
	return nil
}
