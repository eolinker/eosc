package eosc

import (
	"encoding/json"
	"os"
	"reflect"
)

type IDataMarshaller interface {
	Encode(startIndex int) ([]byte, []*os.File, error)
}

func NewBase[T any]() *Base[T] {
	return &Base[T]{}
}

type Base[T any] struct {
	Config *T
	Append map[string]interface{}
}

func (b *Base[T]) UnmarshalJSON(bytes []byte) error {
	cfg := new(T)
	err := json.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}
	appendCfg := make(map[string]interface{})
	err = json.Unmarshal(bytes, &appendCfg)
	if err != nil {
		return err
	}
	typ := reflect.TypeOf(cfg).Elem()
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue // 跳过带有 `json:"-"` 标签的字段
			}
			delete(appendCfg, tag)
		}
	}
	b.Config = cfg
	b.Append = appendCfg
	return nil
}

func (b *Base[T]) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	val := reflect.Indirect(reflect.ValueOf(b.Config)) // 处理指针
	typ := val.Type()

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue // 跳过带有 `json:"-"` 标签的字段
			}
			if _, ok := b.Append[tag]; ok {
				continue // 跳过 b.Append 中的字段
			}
			result[tag] = val.Field(i).Interface()
		}
	}

	// 合并 b.Append 的内容
	for key, v := range b.Append {
		result[key] = v
	}

	return json.Marshal(b.Append)
}

func (b *Base[T]) SetAppend(key string, value interface{}) {
	b.Append[key] = value
}
