package main

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	               "${status_key}":200
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
	//fmt.Println('$', '{', '}')

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

func interfaceDeal(originVal reflect.Value, targetVal reflect.Value, variable map[string]string) error {
	value := originVal.Interface()
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		return mapDeal(v, targetVal, variable, "")
	case reflect.Array, reflect.Slice:
		return arraySet(v, targetVal, variable)
	case reflect.String:
		return stringSet(value.(string), targetVal, variable)
	default:
	}
	return nil
}
