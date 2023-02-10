package dubbo_context

import (
	"github.com/eolinker/eosc/eocontext"
	"time"
)

type IRequestReader interface {
	Header() IHeaderReader
	Service() IServiceReader
	Body() interface{}
	Host() string
	Attachments() map[string]interface{}
	Attachment(string) (interface{}, bool)
}

type IServiceReader interface {
	Path() string
	Interface() string
	Group() string
	Version() string
	Method() string
	Timeout() time.Duration //request timeout
}

type IHeaderReader interface {
	ID() int64
	SerialID() byte
	// Type PackageType
	Type() int
	BodyLen() int
	ResponseStatus() byte
}

type IDubboContext interface {
	eocontext.EoContext
	HeaderReader() IRequestReader // 读取原始请求
	Proxy() IProxy                // 读写转发请求
	Response() IResponse          // 处理返回结果，可读可写
	SendTo(address string, timeout time.Duration) error
}

type IResponse interface {
	ResponseError() error
	SetResponseTime(duration time.Duration)
	ResponseTime() time.Duration
	IBodyGet
	IBodySet
}

type IBodyGet interface {
	GetBody() interface{}
}

type IBodySet interface {
	SetBody(interface{})
}

type IProxy interface {
	Header() IHeaderWriter
	Service() IServiceWriter
	IBodySet
	IBodyGet
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
	SetTimeout(duration time.Duration) //request timeout
}

type IHeaderWriter interface {
	IHeaderReader
	SetID(int64)
	SetSerialID(byte)
	SetType(int)
	SetBodyLen(int)
}
