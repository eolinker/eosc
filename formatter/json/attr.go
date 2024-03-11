package json

import (
	"strconv"
)

var (
	attrOmitempty = "omitempty"
)

var (
	genAttrHandlerMap = map[string]genAttrHandlerFunc{
		attrOmitempty: newOmitEmptyHandler,
	}
)

type IAttrHandler interface {
	Handle(value interface{}, t filedType) (interface{}, bool)
}

type genAttrHandlerFunc = func(args string) IAttrHandler

type omitEmptyHandler struct {
	omitempty bool
}

func newOmitEmptyHandler(args string) IAttrHandler {
	omitempty := true
	if args != "" {
		omitempty, _ = strconv.ParseBool(args)
	}
	return &omitEmptyHandler{omitempty: omitempty}
}

func (o *omitEmptyHandler) Handle(value interface{}, t filedType) (interface{}, bool) {
	switch t {
	case Constants:
	case Variable:
		// 变量
		switch v := value.(type) {
		case string:
			if v == "" && o.omitempty {
				return nil, false
			}
		case int, int64, float32, float64:
			if v == 0 && o.omitempty {
				return nil, false
			}
		}
	case Array:
		v, ok := value.([]interface{})
		if !ok && o.omitempty {
			return nil, false
		}
		if (v == nil || len(v) < 1) && o.omitempty {
			return nil, false
		}
	case Object:
		v, ok := value.(map[string]interface{})
		if !ok && o.omitempty {
			return nil, false
		}
		if (v == nil || len(v) < 1) && o.omitempty {
			return nil, false
		}
	}
	return value, true
}
