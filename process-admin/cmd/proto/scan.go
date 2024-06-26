package proto

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Scan parses bytes `b` to `v` with appropriate type.
//
//nolint:gocyclo
func Scan(b []byte, v interface{}) error {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("apinto: Scan(nil)")
	case *string:
		*v = string(b)
		return nil
	case *[]byte:
		*v = b
		return nil
	case *int:
		var err error
		*v, err = strconv.Atoi(string(b))
		return err
	case *int8:
		n, err := strconv.ParseInt(string(b), 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return nil
	case *int16:
		n, err := strconv.ParseInt(string(b), 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return nil
	case *int32:
		n, err := strconv.ParseInt(string(b), 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return nil
	case *int64:
		n, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *uint:
		n, err := strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return nil
	case *uint8:
		n, err := strconv.ParseUint(string(b), 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := strconv.ParseUint(string(b), 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := strconv.ParseUint(string(b), 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *float32:
		n, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*v = float32(n)
		return err
	case *float64:
		var err error
		*v, err = strconv.ParseFloat(string(b), 64)
		return err
	case *bool:
		*v = len(b) == 1 && b[0] == '1'
		return nil
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, string(b))
		return err
	case *time.Duration:
		n, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = time.Duration(n)
		return nil
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)
	default:
		return fmt.Errorf(
			"apinto: can't unmarshal %T (consider implementing BinaryUnmarshaler)", v)
	}
}

func ScanSlice(data []string, slice interface{}) error {
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
	for i, s := range data {
		elem := next()
		if err := Scan([]byte(s), elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("apinto: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}

	return nil
}

func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value {
	elemType := v.Type().Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		return func() reflect.Value {
			if v.Len() < v.Cap() {
				v.Set(v.Slice(0, v.Len()+1))
				elem := v.Index(v.Len() - 1)
				if elem.IsNil() {
					elem.Set(reflect.New(elemType))
				}
				return elem.Elem()
			}

			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			return elem.Elem()
		}
	}

	zero := reflect.Zero(elemType)
	return func() reflect.Value {
		if v.Len() < v.Cap() {
			v.Set(v.Slice(0, v.Len()+1))
			return v.Index(v.Len() - 1)
		}

		v.Set(reflect.Append(v, zero))
		return v.Index(v.Len() - 1)
	}
}
