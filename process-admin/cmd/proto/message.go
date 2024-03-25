package proto

import (
	"fmt"
	"reflect"
	"strconv"
)

type IMessage interface {
	Array() ([]IMessage, error)
	String() (string, error)
	Int() (int64, error)
	Float() (float64, error)
	Type() ReplyType
	Scan(v any) error
}
type MessageBase []byte

func (m MessageBase) Scan(v any) error {
	return Scan(m, v)
}

type MessageString string

func (m MessageString) Scan(v any) error {
	return Scan([]byte(m), v)
}

func (m MessageString) Array() ([]IMessage, error) {

	return nil, fmt.Errorf(" can't parse array reply: %.100q", m)

}

func (m MessageString) String() (string, error) {
	return string(m), nil
}

func (m MessageString) Int() (int64, error) {
	return strconv.ParseInt(string(m), 10, 64)
}

func (m MessageString) Float() (float64, error) {
	return strconv.ParseFloat(string(m), 64)
}

func (m MessageString) Type() ReplyType {
	return StringReply
}

func (m MessageBase) Array() ([]IMessage, error) {
	switch m[0] {
	case ErrorReply:
		return nil, Error(m[1:])
	default:
		return nil, fmt.Errorf(" can't parse array reply: %.100q", m)
	}
}

func (m MessageBase) String() (string, error) {
	switch m[0] {
	case ErrorReply:
		return "", Error(m[1:])
	case StatusReply:
		return string(m[1:]), nil
	case IntReply:
		return string(m[1:]), nil
	default:
		return "", fmt.Errorf("apinto: can't parse reply=%.100q reading string", m)
	}
}

func (m MessageBase) Int() (int64, error) {
	switch m[0] {
	case ErrorReply:
		return 0, Error(m[1:])
	case IntReply:
		return strconv.ParseInt(string(m[1:]), 10, 64)
	case StatusReply:
		return strconv.ParseInt(string(m[1:]), 10, 64)
	case StringReply:
		return strconv.ParseInt(string(m[1:]), 10, 64)

	default:
		return 0, fmt.Errorf("apinto: can't parse int reply: %.100q", string(m[1:]))
	}
}

func (m MessageBase) Float() (float64, error) {
	switch m[0] {
	case ErrorReply:
		return 0, Error(m[1:])
	case IntReply:
		return strconv.ParseFloat(string(m[1:]), 64)

	case StringReply:
		return strconv.ParseFloat(string(m[1:]), 64)
	case StatusReply:
		return strconv.ParseFloat(string(m[1:]), 64)
	default:
		return 0, fmt.Errorf("apinto: can't parse float reply: %.100q", string(m[1:]))
	}
}

func (m MessageBase) Error() error {
	if m[0] == ErrorReply {
		return Error(m[1:])
	}
	return nil
}

func (m MessageBase) Type() ReplyType {
	return m[0]
}

type MessageArray struct {
	items []IMessage
	err   error
}

func (m *MessageArray) Scan(slice any) error {
	v := reflect.ValueOf(slice)
	if !v.IsValid() {
		return fmt.Errorf("redis: ScanSlice(nil)")
	}
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("redis: ScanSlice(non-pointer %T)", slice)
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("redis: ScanSlice(non-slice %T)", slice)
	}
	next := makeSliceNextElemFunc(v)
	for i, s := range m.items {
		elem := next()
		if err := s.Scan(elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("redis: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}
	return nil
}

func (m *MessageArray) Array() ([]IMessage, error) {
	return m.items, nil
}

func (m *MessageArray) String() (string, error) {
	return "", ErrorCantParseArrayToString
}

func (m *MessageArray) Int() (int64, error) {
	return 0, ErrorCantParseArrayToInt
}

func (m *MessageArray) Float() (float64, error) {
	return 0, ErrorCantParseArrayToFloat
}

func (m *MessageArray) Type() ReplyType {
	return ArrayReply
}
