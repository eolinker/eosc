package json

import (
	json2 "encoding/json"
	"fmt"
	"github.com/eolinker/eosc/formatter"
	"strings"
)

const (
	ROOT            = "fields"
	defaultChildKey = "proxies"
)

var ConfigFormatError = fmt.Errorf("config is not valid")

type filedType string

const (
	Variable  filedType = "variable"
	Constants filedType = "constants"
	Object    filedType = "object"
	Array     filedType = "array"
)

type fieldInfo struct {
	name     string
	cname    string
	t        filedType
	child    map[string]fieldInfo
	childKey string
}

type json struct {
	fields map[string]fieldInfo
}

func (j *json) Format(entry formatter.IEntry) []byte {
	res := make(map[string]interface{})
	if len(j.fields) > 0 {
		res = j.getValue(j.fields, entry)
	}
	b, _ := json2.Marshal(res)
	return b
}

func (j *json) getValue(fields map[string]fieldInfo, entry formatter.IEntry) map[string]interface{} {
	res := make(map[string]interface{})
	for key, info := range fields {
		switch info.t {
		case Constants:
			// 常量
			res[info.cname] = info.name
		case Variable:
			res[info.cname] = entry.Read(key)
		case Array:
			res[info.cname] = j.getArray(info.childKey, info.child, entry)
		case Object:
			res[info.cname] = j.getValue(info.child, entry)
		}
	}
	return res

}

func (j *json) getArray(key string, arr map[string]fieldInfo, entry formatter.IEntry) []interface{} {
	ens := entry.Children(key)
	res := make([]interface{}, 0, len(ens))
	for _, en := range ens {
		res = append(res, j.getValue(arr, en))
	}
	return res
}

func NewFormatter(cfg formatter.Config) (formatter.IFormatter, error) {

	fields, ok := cfg[ROOT]
	if !ok {
		return nil, ConfigFormatError
	}
	data, err := ParseConfig(fields, cfg)
	if err != nil {
		return nil, err
	}
	return &json{fields: data}, nil
}

func ParseConfig(root []string, cfg formatter.Config) (map[string]fieldInfo, error) {
	data := make(map[string]fieldInfo)
	for _, field := range root {
		// 键值,别名,类型
		info := parse(field)
		if info.t != Variable && info.t != Constants {
			c, ok := cfg[info.name]
			if !ok {
				return nil, ConfigFormatError
			}
			child, err := ParseConfig(c, cfg)
			if err != nil {
				return nil, err
			}
			info.child = child
		}
		data[info.name] = info
	}
	return data, nil
}

func parse(filed string) fieldInfo {
	filed = strings.Trim(filed, " ")
	fs := strings.Split(filed, " ")
	l := len(fs)
	key := fs[0]
	name := ""
	if l == 3 && strings.Contains(filed, "as") {
		name = fs[l-1]
	}
	if strings.HasPrefix(key, "$") {
		key = strings.TrimLeft(key, "$")
		name = key
		// 常量
		return fieldInfo{name: key, cname: name, t: Variable}
	}
	if strings.HasPrefix(key, "@") {
		key = strings.TrimLeft(key, "@")
		if strings.Contains(key, "#") {
			if strings.HasSuffix(key, "#") {
				key = strings.TrimRight(key, "#")
				name = key
				// 数组，获取所有字段
				return fieldInfo{name: key, cname: name, t: Array, childKey: defaultChildKey}
			}
			// 数组，获取个别字段
			childKey := strings.Split(key, "#")[1]
			return fieldInfo{name: key, cname: name, t: Array, childKey: childKey}
		}

		name = key
		// 对象
		return fieldInfo{name: key, cname: name, t: Object}
	}
	// 啥都没
	return fieldInfo{name: key, cname: name, t: Constants}
}
