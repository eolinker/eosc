package model

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/professions"
)

type ProfessionConfig = eosc.ProfessionConfig
type ProfessionInfo struct {
	Name   string   `json:"name,omitempty"`
	Label  string   `json:"label,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Driver []string `json:"driver,omitempty"`
}

func TypeProfessionInfo(p *professions.Profession) *ProfessionInfo {
	drivers := make([]string, 0, len(p.Drivers))
	for _, d := range p.Drivers {
		drivers = append(drivers, d.Name)
	}
	return &ProfessionInfo{
		Name:   p.Name,
		Label:  p.Label,
		Desc:   p.Desc,
		Driver: drivers,
	}
}
