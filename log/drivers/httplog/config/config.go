package config

import (
	"encoding/json"
	"github.com/eolinker/goku-standard/common/log"
	"net/http"
)

type encoder interface {
	encode()(string,error)
}
type HeaderItem struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type ConfigEncode struct {

	Method string `json:"method"`
	Url string `json:"url"`
	Headers []HeaderItem `json:"headers"`
	Level string `json:"level"`

}

func (c * ConfigEncode)encode()(string,error)  {

	data, err := json.Marshal(c)
	if err!=nil{
		return  "",err
	}
	return string(data),err
}


type Config struct {

	Method        string
	Url           string
	Headers       http.Header
	Level         log.Level

	HandlerCount int
}

func (c *Config) encode() (string,error) {
	en:= &ConfigEncode{
		Method:c.Method,
		Url:c.Method,
		Level:c.Level.String(),
		Headers:toHeaderItems(c.Headers),

	}
	return en.encode()
}

func toHeaderItems(header http.Header) []HeaderItem {
	headers := make([]HeaderItem, 0, len(header))

	for k:=range header{
		headers = append(headers, HeaderItem{Key:k,
			Value:header.Get(k),
		})
	}
	return headers
}
func toHeader(items []HeaderItem) http.Header {
	header:=make(http.Header)
	for _,item:=range items{
		header.Set(item.Key,item.Value)
	}
	return header
}