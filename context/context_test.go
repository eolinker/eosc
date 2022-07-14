package context

import (
	"fmt"
	"github.com/eolinker/eosc/utils/config"
	"net/http"
)

type HttpContext interface {
	LoadBalance() LoadBalance
	SetLoadBalance(balance LoadBalance)
	DO() DoHandler
	SetDoHandler(handler DoHandler)
	Finish() FinishHandler
	SetFinish(handler FinishHandler)
	Scheme() string
}
type HttpContextAssert interface {
	HttpContext
	Assert(i interface{}) error
}

type ExampleHttpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (e *ExampleHttpContext) Assert(i interface{}) error {
	if v, ok := i.(*HttpContext); ok {
		*v = e
		return nil
	}
	return fmt.Errorf("not suport:%s", config.TypeNameOf(i))
}

func (e *ExampleHttpContext) LoadBalance() LoadBalance {

	panic("implement me")
}

func (e *ExampleHttpContext) SetLoadBalance(balance LoadBalance) {
	panic("implement me")
}

func (e *ExampleHttpContext) DO() DoHandler {
	panic("implement me")
}

func (e *ExampleHttpContext) SetDoHandler(handler DoHandler) {
	panic("implement me")
}

func (e *ExampleHttpContext) Finish() FinishHandler {
	panic("implement me")
}

func (e *ExampleHttpContext) SetFinish(handler FinishHandler) {
	panic("implement me")
}

func (e *ExampleHttpContext) Scheme() string {
	return e.r.URL.Scheme
}

func Example_Context() {

	ctx := ExampleHttpContext{}

	var httpContext HttpContext
	err := ctx.Assert(&httpContext)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	//
}
