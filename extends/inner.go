package extends

import "fmt"

var (
	innerExtender = make(map[string][]RegisterFunc)
)

func AddInnerExtendProject(group, name string, registerFunc ...RegisterFunc) {
	id := fmt.Sprint(group, "-", name)
	innerExtender[id] = append(innerExtender[id], registerFunc...)
}
