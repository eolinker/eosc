package eosc

import (
	"fmt"
	"strings"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ToDriverDetails(config []*DriverConfig) []*DriverDetail {
	rs := make([]*DriverDetail, len(config))
	for i, c := range config {
		rs[i] = ToDriverDetail(c)
	}
	return rs
}
func ToDriverDetail(config *DriverConfig) *DriverDetail {
	group, project, name := readDriverId(config.Id)
	return &DriverDetail{
		Id:     config.Id,
		Driver: config.Name,
		Label:  config.Label,
		Desc:   config.Desc,
		Plugin: &PluginInfo{
			Group:   group,
			Project: project,
			Name:    name,
		},
		Params: config.Params,
	}
}
func readDriverId(id string) (group, project, name string) {
	vs := strings.Split(id, ":")

	if len(vs) > 2 {

		group = vs[0]
		project = vs[1]
		name = vs[2]
		return
	}
	if len(vs) == 2 {
		project = vs[0]
		name = vs[1]
		return
	}
	name = vs[0]
	return
}

func ToWorkerId(name, profession string) string {
	return fmt.Sprintf("%s@%s", name, profession)
}
