package bean

import (
	"log"
)

type NameInterface interface {
	DO()
}
type TestImpl struct {
	name string
}

func (t *TestImpl) DO() {
	log.Println("do by:", t.name)
}

func init() {

	var nameInterface NameInterface
	// 这里一定要要用指针
	Autowired(&nameInterface)

	AddInitializingBeanFunc(func() {
		log.Println("auto wired done")
		nameInterface.DO()
	})

}

func ExampleAutowired() {

	t := &TestImpl{name: "demo bean"}
	// 转换成注入目标的类型
	var i NameInterface = t
	// 这里也要用指针，否则反射识别类型会出问题
	Injection(&i)
	// 检查是否有缺失
	err := Check()
	if err != nil {
		panic(err)
	}
}
