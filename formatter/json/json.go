package json

import (
	json2 "encoding/json"
	"fmt"
	"strings"

	"github.com/eolinker/eosc/formatter"
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

type jsonFormat struct {
	fields map[string]fieldInfo
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
	return &jsonFormat{fields: data}, nil
}

func (j *jsonFormat) Format(entry formatter.IEntry) []byte {
	res := make(map[string]interface{})
	if len(j.fields) > 0 {
		res = j.getValue(j.fields, entry)
	}
	b, _ := json2.Marshal(res)
	return b
}

func (j *jsonFormat) getValue(fields map[string]fieldInfo, entry formatter.IEntry) map[string]interface{} {
	res := make(map[string]interface{})
	for key, info := range fields {
		switch info.t {
		case Constants:
			// 常量
			res[key] = info.name
		case Variable:
			res[key] = entry.Read(info.name)
		case Array:
			res[key] = j.getArray(info.childKey, info.child, entry)
		case Object:
			res[key] = j.getValue(info.child, entry)
		}
	}
	return res

}

func (j *jsonFormat) getArray(key string, arr map[string]fieldInfo, entry formatter.IEntry) []interface{} {
	ens := entry.Children(key)
	res := make([]interface{}, 0, len(ens))
	for _, en := range ens {
		res = append(res, j.getValue(arr, en))
	}
	return res
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
		data[info.cname] = info
	}
	return data, nil
}

func parse(filed string) fieldInfo {
	filed = strings.Trim(filed, " ")
	fs := strings.Split(filed, " ")
	key := fs[0]
	res := fieldInfo{name: key, t: Constants}
	if strings.HasPrefix(key, "$") {
		key = strings.TrimLeft(key, "$")
		// 常量
		res = fieldInfo{name: key, t: Variable}
	}
	if strings.HasPrefix(key, "@") {
		key = strings.TrimLeft(key, "@")
		if strings.Contains(key, "#") {
			if strings.HasSuffix(key, "#") {
				key = strings.TrimRight(key, "#")
				// 数组，获取所有字段
				res = fieldInfo{name: key, t: Array, childKey: defaultChildKey}
			} else {
				// 数组，获取检索字段
				s := strings.Split(key, "#")
				key, childKey := s[0], s[1]
				res = fieldInfo{name: key, t: Array, childKey: childKey}
			}
		} else {
			// 对象
			res = fieldInfo{name: key, t: Object}
		}
	}
	l := len(fs)
	if l == 3 && strings.Contains(filed, "as") {
		res.cname = fs[l-1]
	} else {
		res.cname = res.name
	}
	return res
}
