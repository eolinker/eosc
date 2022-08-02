package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	strArray := []struct {
		origin string
		want   string
	}{
		{
			origin: `{
	   "listen":"${listen}",
	   "method":[
	       "${Get}",
	       "POST"
	   ],
	   "host":[
	
	   ],
	   "rules":[
	       {
	           "location":"/abc",
	           "header":{
	               "abc":"abc"
	           },
	           "query":{
	               "abc":"abc"
	           }
	       }
	   ],
	   "target":"test@service",
	   "disable":"${is_close}",
	   "plugins":{
	       "write":{
	           "disable":false,
	           "config":{
	               "body":"{\"abc\":\"${value}\"}",
	               "body_base64":true,
	               "headers":{
	                   "test":"${listen}"
	               },
	               "match":{
	                   "code":[
	                       "${code}",
	                       300
	                   ]
	               },
	               "status_code":200
	           }
	       }
	   }
	}`,
			want: `{
	   "listen":"8099",
	   "method":[
	       "GET",
	       "POST"
	   ],
	   "host":[
	
	   ],
	   "rules":[
	       {
	           "location":"/abc",
	           "header":{
	               "abc":"abc"
	           },
	           "query":{
	               "abc":"abc"
	           }
	       }
	   ],
	   "target":"test@service",
	   "disable":"true",
	   "plugins":{
	       "write":{
	           "disable":false,
	           "config":{
	               "body":"{\"abc\":\"abc\"}",
	               "body_base64":true,
	               "headers":{
	                   "test":"8099"
	               },
	               "match":{
	                   "code":[
	                       "200",
	                       300
	                   ]
	               },
	               "status_code":200
	           }
	       }
	   }
	}`,
		},
	}

	m := map[string]string{
		"abd":        "replace",
		"abc":        "asdj",
		"listen":     "8099",
		"is_close":   "true",
		"Get":        "\"GET\"",
		"value":      "\"abc\"",
		"code":       "200",
		"status_key": "status_code",
	}
	for _, str := range strArray {
		parse, err := NewParse(reflect.TypeOf(&DriverConfig{}), m)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(str.origin), parse)
		if err != nil {
			fmt.Println("err:", err, "str:", str.origin)
			continue
		}
		marshal, _ := json.Marshal(parse.origin)
		fmt.Println("finally target value:", string(marshal))
	}
}

func interfaceSet(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	value := originVal.Elem()
	switch value.Kind() {
	case reflect.Map:
		return mapSet(value, targetVal, variable, "")
	case reflect.Array, reflect.Slice:
		return arraySet(value, targetVal, variable)
	case reflect.String:
		return stringSet(value, targetVal, variable)
	case reflect.Float64:
		return float64Set(value, targetVal)
	default:
		fmt.Println("interface deal", "kind", value.Kind())
	}
	return nil
}

func float64Set(originVal reflect.Value, targetVal reflect.Value) error {
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}
	switch targetVal.Kind() {
	case reflect.Int:
		value, err := strconv.ParseInt(fmt.Sprintf("%1.0f", originVal.Float()), 10, 64)
		if err != nil {
			return err
		}
		targetVal.SetInt(value)
	case reflect.Float64:
		targetVal.SetFloat(originVal.Float())
	case reflect.String:
		value := fmt.Sprintf("%f", originVal.Float())
		targetVal.SetString(value)
	default:
		return fmt.Errorf("float64 set error:%w %s", ErrorUnsupportedKind, targetVal.Kind())
	}
	return nil
}
