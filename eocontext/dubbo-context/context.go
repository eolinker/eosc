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
	Method() string         //固定值：$invoke
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
	SendTo(address string, timeout time.Duration) error
}

type IProxy interface {
	Header() IHeaderWriter
	Service() IServiceWriter
	SetBody(interface{})
}

type IServiceWriter interface {
	IServiceReader
	SetPath(string)
	SetInterface(string)
	SetGroup(string)
	SetVersion(string)
	SetMethod(string)                  //固定值：$invoke
	SetTimeout(duration time.Duration) //request timeout
}

type IHeaderWriter interface {
	IHeaderReader
	SetID(int64)
	SetSerialID(byte)
	SetType(int)
	SetBodyLen(int)
}
