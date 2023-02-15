package grpc_context

import (
	"time"

	"github.com/eolinker/eosc/eocontext"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc/metadata"
)

type IGrpcContext interface {
	eocontext.EoContext

	//ServerStream() grpc.ServerStream

	// Request 获取原始请求
	Request() IRequest
	// Proxy 获取待转发请求
	Proxy() IRequest
	// Response 获取服务端响应内容
	Response() IResponse

	// EnableTls 是否开启tls认证
	EnableTls(bool)
	// InsecureCertificateVerify 是否跳过证书检查
	InsecureCertificateVerify(bool)
	// Invoke grpc调用
	Invoke(address string, timeout time.Duration) error
	FastFinish() error
}

type IRequest interface {
	Headers() metadata.MD
	Host() string
	Service() string
	SetService(string)
	Method() string
	SetMethod(string)
	FullMethodName() string
	RealIP() string
	ForwardIP() string
	// Message 获取原始请求内容，在grpc协议需要转其他协议时使用
	Message(*desc.MessageDescriptor) *dynamic.Message
}

type IResponse interface {
	Headers() metadata.MD
	Message() *dynamic.Message
	Trailer() metadata.MD
	Write(msg *dynamic.Message)
	SetErr(err error)
	Error() error
	//ClientStream() grpc.ClientStream
}

type ITrailer interface {
	IHeaderReader
	IHeaderWriter
}

type IHeader interface {
	IHeaderReader
	IHeaderWriter
}

type IHeaderReader interface {
	Get(key string) string
	Headers()
}

type IHeaderWriter interface {
	Set(key string, value string)
}
