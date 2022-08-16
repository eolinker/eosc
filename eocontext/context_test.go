package eocontext

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils/config"
	"net/http"
)

type HttpContext interface {
	EoContext
}

type ExampleHttpContext struct {
	w http.ResponseWriter
	r *http.Request

	complete CompleteHandler
	finish   FinishHandler
}

func (e *ExampleHttpContext) SetLabel(name, value string) {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) GetLabel(name string) string {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) Labels() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) RequestId() string {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) Context() context.Context {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) Value(key interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (e *ExampleHttpContext) WithValue(key, val interface{}) {
	//TODO implement me
	panic("implement me")
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

	var ctx EoContext = &ExampleHttpContext{}

	var httpContext HttpContext
	err := ctx.Assert(&httpContext)
	if err != nil {
		log.Debug(err)
		return
	}
	// Output:
	//
}
