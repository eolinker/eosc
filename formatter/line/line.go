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

type Container struct {
	left  rune
	right rune
}

type Line struct {
	cfg formatter.Config
}

func (l *Line) Format(entry formatter.IEntry) []byte {
	fields, ok := l.cfg["fields"]
	if !ok {
		return []byte("")
	}
	values := l.recursionField(fields, entry, 0)
	data := strings.Join(values, separators[0])
	return []byte(data)
}

func (l *Line) recursionField(fields []string, entry formatter.IEntry, level int) []string {
	data := make([]string, len(fields))
	if separatorLen <= level {
		return nil
	}
	cta := containers[level]
	left := string(cta.left)
	right := string(cta.right)
	for i, name := range fields {
		value, ok := entry.Read(name)
		if !ok {
			if strings.HasPrefix(name, "@") {
				value = ""
				if strings.HasSuffix(name, "#") {
					// array
				} else {
					n := strings.TrimPrefix(name, "@")
					fs, ok := l.cfg[n]
					if ok {
						result := l.recursionField(fs, entry, level+1)
						if separatorLen < level+1 {
							continue
						}
						value = strings.Join(result, separators[level+1])
					}
				}
			} else if strings.HasPrefix(name, "$") {
				value = ""
			} else {
				value = name
			}
		}
		data[i] = left + value + right
	}
	return data
}
