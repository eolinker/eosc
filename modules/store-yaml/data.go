package store

import (
	"encoding/json"
	"github.com/eolinker/eosc"
)

var _ eosc.IData = (*Router)(nil)

type Cert struct {
	Key string `json:"key" yaml:"key"`
	Crt string `json:"crt" yaml:"crt"`
}


type Rule struct {
	Target string `json:"target" yaml:"target"`
	RemoteIp string `json:"ip" yaml:"ip"`
	Host string `json:"host" yaml:"host"`
	Location string `json:"location" yaml:"location"`
	Header map[string]string `json:"header" yaml:"header"`
	Query map[string]string `json:"query" yaml:"query"`
	Cookie map[string]string `json:"cookie" yaml:"cookie"`
}


type Router struct {
	Name string `json:"name" yaml:"name"`
	Driver string `json:"driver" yaml:"driver"`
	Listen int `json:"listen" yaml:"listen"`
	Host []string `json:"host" yaml:"host"`
	Cert []Cert `json:"cert" yaml:"cert"`
	Rules []Rule `json:"rules" yaml:"rule"`
}

func (r *Router) Marshal() ([]byte, error) {
	d,err:=json.Marshal(r)
	return d,err
}

func (r *Router) UnMarshal(v interface{}) error {

	d,err:=json.Marshal(r)
	if err!= nil{
		return err
	}
	return json.Unmarshal(d,v)
}

type Service struct {
	Name string `json:"name" yaml:"name"`
	Driver string `json:"driver" yaml:"driver"`
}

type Config struct {
	Include []string `json:"include" yaml:"router"`
	Router []Router `json:"router" yaml:"router"`
	Service []Service `json:"service" yaml:"service"`
}


