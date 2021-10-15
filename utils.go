package eosc

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetRealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.RemoteAddr
	}
	return realIP
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

func ToWorkerId(value, profession string) (string, bool) {
	value = strings.ToLower(value)
	index := strings.Index(value, "@")
	if index < 0 {
		return fmt.Sprintf("%s@%s", value, profession), true
	}
	if profession != value[index+1:] {
		return "", false
	}
	return value, true
}
