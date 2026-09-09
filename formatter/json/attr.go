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
	case Number:
		// 统一转成float64
		switch v := value.(type) {
		case int:
			if v == 0 && o.omitempty {
				return nil, false
			}
			return float64(v), true
		case int64:
			if v == 0 && o.omitempty {
				return nil, false
			}
			return float64(v), true
		case float32:
			if v == 0 && o.omitempty {
				return nil, false
			}
			return float64(v), true
		case float64:
			if v == 0 && o.omitempty {
				return nil, false
			}
			return v, true
		case string:
			if v == "" && o.omitempty {
				return nil, false
			}
			f, _ := strconv.ParseFloat(v, 64)

			return f, true
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
