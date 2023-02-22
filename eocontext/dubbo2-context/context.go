package dubbo2_context

import (
	"github.com/eolinker/eosc/eocontext"
	"time"
)

type IRequestReader interface {
	Service() IServiceReader
	Body() interface{}
	Host() string
	Attachments() map[string]interface{}
	Attachment(string) (interface{}, bool)
	RemoteIP() string
}

type IServiceReader interface {
	Path() string
	Interface() string
	Group() string
	Version() string
	Method() string
}

type IDubbo2Context interface {
	eocontext.EoContext
	HeaderReader() IRequestReader // 读取原始请求
	Proxy() IProxy                // 读写转发请求
	Response() IResponse          // 处理返回结果，可读可写
	Invoke(address string, timeout time.Duration) error
}

type IResponse interface {
	ResponseError() error
	SetResponseTime(duration time.Duration)
	ResponseTime() time.Duration
	GetBody() interface{}
	SetBody(interface{})
}

type IProxy interface {
	Service() IServiceWriter
	SetParam(*Dubbo2ParamBody)
	GetParam() *Dubbo2ParamBody
	SetAttachment(string, interface{})
	Attachments() map[string]interface{}
}

type IServiceWriter interface {
	IServiceReader
	SetPath(string)
	SetInterface(string)
	SetGroup(string)
	SetVersion(string)
	SetMethod(string)
}
