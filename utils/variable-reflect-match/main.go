package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var typeMap = map[string]reflect.Type{
	"proxy": reflect.TypeOf(new(Config)),
	"write": reflect.TypeOf(new(RewriteConfig)),
}

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
            "disable":true,
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
        },
        "proxy":{
            "disable":true,
            "config":{
                "scheme":"${scheme}",
                "uri":"/abc",
                "regex_uri":[
                    "/${code}/test",
                    "/test/$1"
                ],
                "host":"demo.apinto.com",
                "headers":{

                }
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
            "disable":true,
            "config":{
                "body":"{\"abc\":\"abc\"}",
                "body_base64":true,
                "headers":{
                    "test":"8099"
                },
                "match":{
                    "code":[
                        200,
                        300
                    ]
                },
                "status_code":200
            }
        },
        "proxy":{
            "disable":true,
            "config":{
                "scheme":"http",
                "uri":"/abc",
                "regex_uri":[
                    "/200/test",
                    "/test/$1"
                ],
                "host":"demo.apinto.com",
                "headers":{

                }
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
		"scheme":     "200",
	}
	parse, err := NewParse(m)
	if err != nil {
		panic(err)
	}

	for _, str := range strArray {
		value, err := parse.Unmarshal([]byte(str.origin), reflect.TypeOf(new(DriverConfig)))
		if err != nil {
			panic(err)
		}
		marshal, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		fmt.Println("finally target value:", string(marshal))
	}
}
