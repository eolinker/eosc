package config

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/utils/schema"
	_ "github.com/stretchr/testify/assert"
	"reflect"
)

func Example() {
	type MyConfig struct {
		Id     string    `json:"id" require:"" readonly:"true"`
		Target RequireId `json:"target" skill:"service.service.IService"`
	}
	sc, err := schema.Generate(reflect.TypeOf(MyConfig{}), nil)
	if err != nil {
		log.Debug(err)
		return
	}
	data, _ := json.MarshalIndent(sc, "", "\t")
	log.Debug(string(data))
	//output: ""
}
