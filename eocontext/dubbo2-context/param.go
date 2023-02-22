package dubbo2_context

import hessian "github.com/apache/dubbo-go-hessian2"

type Dubbo2ParamBody struct {
	TypesList  []string
	ValuesList []hessian.Object
}
