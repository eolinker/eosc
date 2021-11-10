package http

//Checker 路由指标检查器接口
type Checker interface {
	Check(v string, has bool) bool
	Key() string
	CheckType() CheckType
	Value() string
}
