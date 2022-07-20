package context

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils/config"
	"net/http"
)

type HttpContext interface {
	Context
}

type ExampleHttpContext struct {
	w http.ResponseWriter
	r *http.Request

	complete CompleteHandler
	finish   FinishHandler
}

func (e *ExampleHttpContext) Complete() CompleteHandler {
	return e.complete
}

func (e *ExampleHttpContext) SetCompleteHandler(handler CompleteHandler) {
	e.complete = handler
}

func (e *ExampleHttpContext) Assert(i interface{}) error {
	if v, ok := i.(*HttpContext); ok {
		*v = e
		return nil
	}
	return fmt.Errorf("not suport:%s", config.TypeNameOf(i))
}

func (e *ExampleHttpContext) Finish() FinishHandler {
	return e.finish
}

func (e *ExampleHttpContext) SetFinish(handler FinishHandler) {
	e.finish = handler
}

func (e *ExampleHttpContext) Scheme() string {
	return e.r.URL.Scheme
}

func Example_Context() {

	var ctx Context = &ExampleHttpContext{}

	var httpContext HttpContext
	err := ctx.Assert(&httpContext)
	if err != nil {
		log.Debug(err)
		return
	}
	// Output:
	//
}
