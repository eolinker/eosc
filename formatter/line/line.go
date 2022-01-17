package line

import (
	"encoding/base64"
	"strings"

	"github.com/eolinker/eosc"
)

var (
	containers = []Container{
		{
			left:  '"',
			right: '"',
		},
		{
			left:  '[',
			right: ']',
		},
		{
			left:  '<',
			right: '>',
		},
	}
	objFields                = toSet([]string{"request_body", "proxy_body", "response", "response_body"})
	separators, separatorLen = toArr("\t ,|")

	containerLen = len(containers)
)

const (
	constant = iota
	variable
	object
	arr
)

func toArr(v string) ([]string, int) {
	ls := make([]string, len(v))
	for i := 0; i < len(v); i++ {
		ls[i] = v[i : i+1]
	}
	return ls, len(ls)
}
func toSet(arr []string) map[string]bool {
	set := make(map[string]bool)
	for _, k := range arr {
		set[k] = true
	}
	return set
}

type Container struct {
	left  rune
	right rune
}

type Line struct {
	executors map[string][]*executor
}

type executor struct {
	fieldType int
	key       string
	child     string
}

func NewLine(cfg eosc.FormatterConfig) (*Line, error) {
	executors := make(map[string][]*executor)
	for key, strArr := range cfg {
		extList := make([]*executor, len(strArr))

		for i, str := range strArr {
			ext := new(executor)
			//切割，除去as等多余字符串
			newStr := strings.Split(str, " ")[0]
			//对str进行处理，分类四种类型
			if strings.HasPrefix(newStr, "$") {
				ext.fieldType = variable
				ext.key = strings.TrimPrefix(newStr, "$")
			} else if strings.HasPrefix(newStr, "@") {
				newStr = strings.TrimPrefix(newStr, "@")
				if idx := strings.Index(newStr, "#"); idx != -1 {
					ext.fieldType = arr
					ext.key = newStr[:idx]
					ext.child = newStr[idx+1:]
				} else {
					ext.fieldType = object
					ext.key = newStr
				}

			} else {
				ext.fieldType = constant
				ext.key = newStr
			}

			extList[i] = ext
		}

		executors[key] = extList
	}

	return &Line{executors: executors}, nil
}

func (l *Line) Format(entry eosc.IEntry) []byte {
	fields, ok := l.executors["fields"]
	if !ok || len(fields) == 0 {
		return []byte("")
	}

	values := l.recursionField(fields, entry, 0)
	data := strings.Join(values, separators[0])

	return []byte(data)
}

func (l *Line) recursionField(fields []*executor, entry eosc.IEntry, level int) []string {
	data := make([]string, len(fields))
	if separatorLen <= level {
		return []string{}
	}

	var left, right string
	if containerLen > level {
		cta := containers[level]
		left = string(cta.left)
		right = string(cta.right)
	}
	nextLevel := level + 1
	arrayLevel := level + 2

	for i, ext := range fields {

		switch ext.fieldType {
		case constant:
			data[i] = ext.key
		case variable:
			value := entry.Read(ext.key)
			if objFields[ext.key] {
				value = base64.StdEncoding.EncodeToString([]byte(value))
			}
			if value == "" {
				value = "-"
			}
			data[i] = value
		case object:
			fs, ok := l.executors[ext.key]
			value := "-"
			if ok && separatorLen > nextLevel {

				result := l.recursionField(fs, entry, nextLevel)
				value = left + strings.Join(result, separators[nextLevel]) + right
			}
			data[i] = value
		case arr:
			value := "-"
			fs, ok := l.executors[ext.key]
			if ok && separatorLen > level+1 {
				entryList := entry.Children(ext.child)
				results := make([]string, len(entryList))

				var arrLeft, arrRight string
				if containerLen > level+1 {
					cta := containers[level+1]
					arrLeft = string(cta.left)
					arrRight = string(cta.right)
				}
				for idx, e := range entryList {
					if separatorLen > arrayLevel {
						result := l.recursionField(fs, e, arrayLevel)
						results[idx] = arrLeft + strings.Join(result, separators[arrayLevel]) + arrRight
						continue
					}
					results[idx] = "-"
				}
				value = left + strings.Join(results, separators[nextLevel]) + right
			}
			data[i] = value
		}

	}
	return data
}
