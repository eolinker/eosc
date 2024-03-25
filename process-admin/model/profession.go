package model

import "github.com/eolinker/eosc"

type ProfessionConfig = eosc.ProfessionConfig
type ProfessionInfo struct {
	Name   string   `json:"name,omitempty"`
	Label  string   `json:"label,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Driver []string `json:"driver,omitempty"`
}
