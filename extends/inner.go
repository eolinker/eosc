package extends

import (
	"fmt"
	"sync"

	"github.com/eolinker/eosc"
)

var (
	innerLock     sync.Mutex
	innerExtender = make(map[string]map[string][]RegisterFunc)
	projectCount  = 0
)

func AddInnerExtendProject(group, project string, registerFunc ...RegisterFunc) {
	innerLock.Lock()
	defer innerLock.Unlock()
	projects, has := innerExtender[group]
	if !has {
		projects = make(map[string][]RegisterFunc)
		innerExtender[group] = projects
	}
	if hs, has := projects[project]; !has {
		projectCount++
		hs = make([]RegisterFunc, 0, 10)
		projects[project] = append(hs, registerFunc...)

	} else {
		projects[project] = append(hs, registerFunc...)

	}
}

//lookInner 查看内置插件
func lookInner(group, project string) ([]RegisterFunc, bool) {
	projects, has := innerExtender[group]
	if has {
		rfs, ok := projects[project]
		return rfs, ok
	}
	return nil, false
}

//LoadInner 加载内置插件
func LoadInner(register eosc.IExtenderDriverRegister) {
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
func GetInners() []string {
	innerLock.Lock()
	defer innerLock.Unlock()
	rs := make([]string, 0, projectCount)
	for group, projects := range innerExtender {
		for project := range projects {

			rs = append(rs, fmt.Sprint(group, ":", project))
		}
	}
	return rs
}
