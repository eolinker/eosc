package proto

import (
	"fmt"
	"reflect"
	"strconv"
)

type IMessage interface {
	Array() (ArrayMessage, error)
	String() (string, error)
	Int() (int64, error)
	Float() (float64, error)
	Type() ReplyType
	Scan(v any) error
}
type ArrayMessage []IMessage

func (arm ArrayMessage) Scan(vs ...any) error {
	max := len(arm)
	if len(vs) < max {
		max = len(vs)
	} else if len(vs) > max {
		return fmt.Errorf("too many arguments")
	}
	for i := 0; i < max; i++ {
		if err := arm[i].Scan(vs[i]); err != nil {
			return err
		}
	}
	return nil
}

type MessageBase []byte

func (m MessageBase) Scan(v any) error {
	return Scan(m[1:], v)
}

type MessageString string

func (m MessageString) Scan(v any) error {
	return Scan([]byte(m), v)
}

func (m MessageString) Array() (ArrayMessage, error) {

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

func (m MessageBase) Array() (ArrayMessage, error) {
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

//
//func (m MessageBase) Error() error {
//	if m[0] == ErrorReply {
//		return Error(m[1:])
//	}
//	return nil
//}

func (m MessageBase) Type() ReplyType {
	return m[0]
}

type MessageArray struct {
	items ArrayMessage
	err   error
}

func (m *MessageArray) Scan(slice any) error {
	v := reflect.ValueOf(slice)
	if !v.IsValid() {
		return fmt.Errorf("apinto: ScanSlice(nil)")
	}
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("apinto: ScanSlice(non-pointer %T)", slice)
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("apinto: ScanSlice(non-slice %T)", slice)
	}
	next := makeSliceNextElemFunc(v)
	for i, s := range m.items {
		elem := next()
		if err := s.Scan(elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("apinto: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}
	return nil
}

func (m *MessageArray) Array() (ArrayMessage, error) {
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
