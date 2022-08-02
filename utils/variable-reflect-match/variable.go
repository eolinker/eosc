package main

import (
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

func NewBuilder(str string) *Builder {
	return &Builder{str: str}
}

type Builder struct {
	str string
}

func (b *Builder) Replace(variableMap map[string]string) (string, bool) {
	strBuilder := strings.Builder{}
	varBuilder := strings.Builder{}
	status := CurrentStatus
	startIndex := 0
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
			v, ok := variableMap[varBuilder.String()]
			if !ok {
				// 变量不存在，报错
				return "", false
			}
			strBuilder.WriteString(v)
			varBuilder.Reset()
			startIndex = i + 1
			status = CurrentStatus
		}

	}
	if status == InputStatus {
		strBuilder.WriteString(b.str[startIndex:])
	}
	return strBuilder.String(), true
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
