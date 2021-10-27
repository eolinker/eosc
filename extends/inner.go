package extends

import (
	"sync"

	"github.com/eolinker/eosc"
)

var (
	innerLock     sync.Mutex
	innerExtender = make(map[string]map[string][]RegisterFunc)
)

func AddInnerExtendProject(group, project string, registerFunc ...RegisterFunc) {
	innerLock.Lock()
	defer innerLock.Unlock()
	projects, has := innerExtender[group]
	if !has {
		projects = make(map[string][]RegisterFunc)
		innerExtender[group] = projects
	}
	projects[project] = append(projects[project], registerFunc...)
}
func lookInner(group, project string) ([]RegisterFunc, bool) {

	projects, has := innerExtender[group]
	if has {
		rfs, ok := projects[project]
		return rfs, ok
	}
	return nil, false
}
func LoadInner(register eosc.IExtenderRegister) {
	innerLock.Lock()
	defer innerLock.Unlock()

	for group, projects := range innerExtender {
		for project, funs := range projects {
			reg := NewExtenderRegister(group, project)
			for _, fun := range funs {
				fun(reg)
			}
			reg.RegisterTo(register)
		}
	}
}
