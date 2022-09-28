package bean

import "github.com/eolinker/eosc/log"

var (
	//Default 默认的bean
	Default = NewContainer()
)

// AutowiredManager 声明需要注入的接口变量
// 如果目标接口还没有注入实例，会在注入后给将接口实例赋值给指针
// 如果目标接口类型已经被注入，会立刻获得有效的接口实例
func Autowired(is ...interface{}) {
	for _, i := range is {
		Default.Autowired(i)
	}
}

// Injection 注入一个构造好的实例
/*
如果注入多个相同接口类型，后注入的实例会覆盖先注入的实例
*/
func Injection(i interface{}) {
	Default.Injection(i)
}

// InjectionDefault 注入一个构造好的默认实例
/**
如果注入多个相同接口类型，后注入的实例会覆盖先注入的实例，default实例不会覆盖普通实例
*/
func InjectionDefault(i interface{}) {
	Default.InjectionDefault(i)
}

//Check 对默认的bean容器执行检查， 如果所有Autowired需求都被满足，返回nil，否则返回与缺失实例有关都error
func Check() error {
	err := Default.Check()
	if err != nil {
		log.Debug("bean.check:", err)
	}
	return err
}

//AddInitializingBean 注册完成回调接口， 执行check成功后会调用回调接口
func AddInitializingBean(handler InitializingBeanHandler) {
	Default.AddInitializingBean(handler)
}

//AddInitializingBeanFunc 注册完成回调方法， 执行check成功后会调用回调方法
func AddInitializingBeanFunc(handler func()) {
	Default.AddInitializingBeanFunc(handler)
}
