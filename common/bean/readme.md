# 自动注入


demo
```go

package main

import (
	"fmt"
	"github.com/eolinker/goku-standard/common/bean"
)

type AutowiredTester interface {
	Name()string
}
type AutowiredTester1 struct {

}

func (a *AutowiredTester1) Name() string {
	return "AutowiredTester1"
}
var(
  tester AutowiredTester
)
func init() {
	// 依赖注入，这里必须用指针
	bean.Autowired(&tester)
}
func main() {

	var t1 AutowiredTester = new(AutowiredTester1)

	// 注入 AutowiredTester， 这里必须用指针
	bean.Injection(&t1)

	bean.Check()// 检查是否完成了完整注入

	fmt.Println(tester.Name())
}

```