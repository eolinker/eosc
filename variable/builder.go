package variable

import (
	"fmt"
	"strings"
)

const (
	dollarSign = 36  // $
	leftSign   = 123 // {
	rightSign  = 125 // }
)

const (

	//CurrentStatus 普通状态
	CurrentStatus = iota

	//ReadyStatus 预备状态
	ReadyStatus

	//InputStatus 输入状态
	InputStatus

	//EndInputStatus 结束输入状态
	EndInputStatus
)

func NewBuilder(str string, separator, suffix string) *Builder {
	return &Builder{str: str, separator: separator, defaultSuffix: suffix}
}

type Builder struct {
	str           string
	separator     string
	defaultSuffix string
}

func (b *Builder) Replace(variables map[string]string) (string, []string, bool) {
	strBuilder := strings.Builder{}
	varBuilder := strings.Builder{}
	status := CurrentStatus
	startIndex := 0
	useVariable := make([]string, 0, len(variables))
	for i, s := range b.str {
		oldStatus := status
		status = toggleStatus(status, s)
		switch status {
		case CurrentStatus:
			if oldStatus == ReadyStatus {
				strBuilder.WriteString(b.str[startIndex : i+1])
				startIndex = i + 1
				continue
			}
			strBuilder.WriteRune(s)
		case ReadyStatus:
			startIndex = i
		case InputStatus:
			if oldStatus == ReadyStatus {
				// 刚切换状态，忽略此时的字符
				continue
			}
			varBuilder.WriteRune(s)
		case EndInputStatus:
			tmp := varBuilder.String()
			v, ok := variables[tmp]
			if !ok {
				// 变量不存在，报错
				return "", nil, false
			}
			if b.separator != "" && b.defaultSuffix != "" {
				index := strings.Index(tmp, b.separator)
				if index == -1 {
					tmp = fmt.Sprintf("%s%s%s", tmp, b.separator, b.defaultSuffix)
				}
			}
			useVariable = append(useVariable, tmp)
			strBuilder.WriteString(v)
			varBuilder.Reset()
			startIndex = i + 1
			status = CurrentStatus
		}

	}
	if status == InputStatus {
		strBuilder.WriteString(b.str[startIndex:])
	}

	return strBuilder.String(), useVariable, true
}

func toggleStatus(status int, c rune) int {
	switch status {
	case CurrentStatus, EndInputStatus:
		if c == dollarSign {
			return ReadyStatus
		}
		return CurrentStatus
	case ReadyStatus:
		if c == leftSign {
			return InputStatus
		}
		return CurrentStatus
	case InputStatus:
		if c == rightSign {
			return EndInputStatus
		}
		return InputStatus
	}
	return status
}
