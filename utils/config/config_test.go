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
		fmt.Println(err)
		return
	}
	data, _ := json.MarshalIndent(sc, "", "\t")
	fmt.Println(string(data))
	//output: ""
}
