package json

import (
	json2 "encoding/json"
	"fmt"
	"strings"

	"github.com/eolinker/eosc"
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
	name  string
	cname string

	t            filedType
	child        []fieldInfo
	childKey     string
	attrHandlers []IAttrHandler
}

type jsonFormat struct {
	fields []fieldInfo
	ctRs   []contentResize
}

type contentResize struct {
	Size   int    `json:"size"`
	Suffix string `json:"suffix"`
}

func NewFormatter(cfg eosc.FormatterConfig, ctRs []contentResize) (eosc.IFormatter, error) {

	fields, ok := cfg[ROOT]
	if !ok {
		return nil, ConfigFormatError
	}
	data, err := ParseConfig(fields, cfg)
	if err != nil {
		return nil, err
	}
	return &jsonFormat{fields: data, ctRs: ctRs}, nil
}

func (j *jsonFormat) Format(entry eosc.IEntry) []byte {
	res := make(map[string]interface{})
	if len(j.fields) > 0 {
		res = j.getValue(j.fields, entry)
	}
	b, _ := json2.Marshal(res)
	return b
}

func (j *jsonFormat) getValue(fields []fieldInfo, entry eosc.IEntry) map[string]interface{} {
	res := make(map[string]interface{})
	tmp := make(map[string]interface{})

	for _, info := range fields {
		ok := true
		var value interface{}
		switch info.t {
		case Constants:
			// 常量
			value = info.name
		case Variable:
			var has bool
			value, has = tmp[info.name]
			if !has {
				value = entry.Read(info.name)
				for _, c := range j.ctRs {
					if strings.HasSuffix(info.name, c.Suffix) {
						if c.Size > 0 && len(value.(string)) > c.Size<<20 {
							value = value.(string)[:c.Size]
							tmp[fmt.Sprintf("%s_complete", info.name)] = 0
						} else {
							tmp[fmt.Sprintf("%s_complete", info.name)] = 1
						}
					}
				}
			}
		case Array:
			value = j.getArray(info.childKey, info.child, entry)
		case Object:
			value = j.getValue(info.child, entry)
		}

		for _, handler := range info.attrHandlers {
			value, ok = handler.Handle(value, info.t)
		}
		if ok {
			res[info.cname] = value
		}
	}
	return res

}

func (j *jsonFormat) getArray(key string, arr []fieldInfo, entry eosc.IEntry) []interface{} {
	ens := entry.Children(key)
	res := make([]interface{}, 0, len(ens))
	for _, en := range ens {
		res = append(res, j.getValue(arr, en))
	}
	return res
}

func ParseConfig(root []string, cfg eosc.FormatterConfig) ([]fieldInfo, error) {
	data := make([]fieldInfo, 0, len(root))
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
		data = append(data, info)
	}
	return data, nil
}

func parse(field string) fieldInfo {
	field = strings.TrimSpace(field)
	// 以分号拆分属性列
	fields := strings.Split(field, ";")
	fs := strings.Split(fields[0], " ")
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
	if len(fields) > 1 {
		res.attrHandlers = parseAttr(fields[1:])
	}
	l := len(fs)
	if l == 3 && strings.Contains(field, "as") {
		res.cname = fs[l-1]
	} else {
		res.cname = res.name
	}
	return res
}

func parseAttr(attrs []string) []IAttrHandler {
	if len(attrs) == 0 {
		return nil
	}
	handlers := make([]IAttrHandler, 0, len(attrs))
	for _, attr := range attrs {
		attr = strings.TrimSpace(attr)
		as := strings.SplitN(attr, ":", 1)
		if v, ok := genAttrHandlerMap[as[0]]; ok {
			arg := ""
			if len(as) > 1 {
				arg = as[1]
			}
			handlers = append(handlers, v(arg))
		}
	}
	return handlers
}
