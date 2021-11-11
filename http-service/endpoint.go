package http_service

//IRouterEndpoint 实现了返回路由规则信息方法的接口，如返回location、Host、IHeader、Query
type IEndpoint interface {
	Location() (Checker, bool)
	Header(name string) (Checker, bool)
	Query(name string) (Checker, bool)
	Headers() []string
	Queries() []string
}
