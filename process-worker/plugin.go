package process_worker

import (
	"fmt"
	"strings"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/extends"
)

func loadPluginEnv(settings map[string]string) ExtenderRegister {
	register := eosc.NewExtenderRegister()
	extends.LoadInner(register)
	inners := extends.GetInners()
	for group, project := range inners {
		delete(settings, toId(group, project))
	}
	for id, version := range settings {
		group, project, ok := readProject(id)
		if ok {
			extenderProject, err := extends.ReadExtenderProject(group, project, version)
			if err != nil {
				continue
			}
			extenderProject.RegisterTo(register)
		}
	}
	return register
}

func readProject(id string) (string, string, bool) {
	i := strings.Index(id, ":")
	if i < 0 {
		return "", "", false
	}
	return id[:i], id[i+1:], true
}
func toId(group, project string) string {
	return fmt.Sprint(group, ":", project)
}
