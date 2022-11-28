package extends

import (
	"fmt"
	"strings"

	"github.com/eolinker/eosc"
)

func InitRegister() IExtenderRegister {
	register := eosc.NewExtenderRegister()
	LoadInner(register)
	return register
}
func LoadProject(group, project, version string, register IExtenderRegister) {
	extenderProject, err := ReadExtenderProject(group, project, version)
	if err != nil {
		return
	}
	extenderProject.RegisterTo(register)
}
func LoadPlugins(plugins map[string]string, register IExtenderRegister) {
	for id, version := range plugins {
		group, project, ok := ReadProject(id)
		if ok {
			if IsInner(group, project) {
				continue
			}
			LoadProject(group, project, version, register)
		}
	}
}

//func LoadPluginEnv(settings map[string]string) IExtenderRegister {
//	register := eosc.NewExtenderRegister()
//	LoadInner(register)
//	for id, version := range settings {
//
//		group, project, ok := ReadProject(id)
//		if ok {
//			if IsInner(group, project) {
//				continue
//			}
//			extenderProject, err := ReadExtenderProject(group, project, version)
//			if err != nil {
//				continue
//			}
//			extenderProject.RegisterTo(register)
//		}
//	}
//	return register
//}

func ReadProject(id string) (string, string, bool) {
	i := strings.Index(id, ":")
	if i < 0 {
		return "", "", false
	}
	return id[:i], id[i+1:], true
}

func toId(group, project string) string {
	return fmt.Sprint(group, ":", project)
}
