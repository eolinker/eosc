package line

import (
	"strings"

	"github.com/eolinker/eosc/formatter"
)

var separators = []string{
	"\t",
	" ",
	",",
	":",
}

var containers = []Container{
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

var (
	separatorLen = len(separators)
	containerLen = len(containers)
)

const (
	constant = iota
	variable
	object
	arr
)

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
}

func NewLine(cfg formatter.Config) *Line {
	executors := make(map[string][]*executor)
	for key, strArr := range cfg {
		extList := make([]*executor, 0, len(strArr))

		for i, str := range strArr {
			ext := new(executor)
			//切割，除去as等多余字符串
			newStr := strings.Split(str, " ")[0]
			//对str进行处理，分类四种类型
			if strings.HasPrefix(newStr, "$") {
				ext.fieldType = variable
				ext.key = strings.Trim(newStr, "$")
			} else if strings.HasPrefix(newStr, "@") {
				newStr = strings.TrimPrefix(newStr, "@")
				if strings.HasSuffix(newStr, "#") {
					newStr = strings.TrimSuffix(newStr, "#")
					ext.fieldType = arr
					ext.key = newStr
					continue
				}

				ext.fieldType = object
				ext.key = newStr
			} else {
				ext.fieldType = constant
				ext.key = newStr
			}

			extList[i] = ext
		}

		executors[key] = extList
	}

	return &Line{executors: executors}
}

func (l *Line) Format(entry formatter.IEntry) []byte {
	fields, ok := l.executors["fields"]
	if !ok {
		return []byte("")
	}

	values := l.recursionField(fields, entry, 0)
	data := strings.Join(values, separators[0])

	return []byte(data)
}

func (l *Line) recursionField(fields []*executor, entry formatter.IEntry, level int) []string {
	data := make([]string, len(fields))
	if separatorLen <= level {
		return nil
	}

	var left, right string
	if containerLen > level {
		cta := containers[level]
		left = string(cta.left)
		right = string(cta.right)
	}

	for i, ext := range fields {

		switch ext.fieldType {
		case constant:
			data[i] = ext.key
		case variable:
			value, _ := entry.Read(ext.key)
			data[i] = value
		case object:
			fs, ok := l.executors[ext.key]
			var value string
			if ok {
				result := l.recursionField(fs, entry, level+1)
				if separatorLen <= level+1 {
					continue
				}

				value = left + strings.Join(result, separators[level+1]) + right
			}
			data[i] = value
		case arr:

		}

	}
	return data
}
