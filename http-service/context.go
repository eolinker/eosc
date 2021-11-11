package http_service

import (
	"context"
	"net/http"
	"net/textproto"
	"net/url"
)

type IHttpContext interface {
	Context() context.Context
	Value(key interface{}) interface{}
	WithValue(key, val interface{})

	ResponseWriter // 处理返回
	RequestId() string
	Request() IRequestReader
	Proxy() IRequest                // 请求信息，包含原始请求数据以及被更高优先级处理过的结果
	ProxyResponse() IResponseReader // 转发后返回的结果
	Finish()
	// 一个key只可以set一次，重复set报错
	SetStoreValue(key string, value interface{}) error
	GetStoreValue(key string) (interface{}, bool)
}

type IHeaderReader interface {
	GetHeader(name string) string
	// 返回所有header，返回值为一个副本，对他的修改不会生效
	Headers() http.Header
}
type IHeaderWriter interface {
	SetHeader(key, value string)
	AddHeader(key, value string)
	DelHeader(key string)
}
type IHeader interface {
	IHeaderReader
	IHeaderWriter
}

//type ICookieReader interface {
//	Cookie(name string) (*http.Cookie, error)
//	Cookies() []*http.Cookie
//}
type ICookieWriter interface {
	AddCookie(c *http.Cookie)
}
type IBodyGet interface {
	GetBody() []byte
}

type IBodySet interface {
	SetBody([]byte)
}
type FileHeader struct {
	FileName string
	Header   textproto.MIMEHeader
	Data     []byte
}
type IBodyDataReader interface {
	//protocol() RequestType
	ContentType() string
	//content-Type = application/x-www-form-urlencoded 或 multipart/form-data，与原生request.Form不同，这里不包括 query 参数
	BodyForm() (url.Values, error)
	//content-Type = multipart/form-data 时有效
	Files() (map[string]*FileHeader, error)
	GetForm(key string) string
	GetFile(key string) (file *FileHeader, has bool)
	RawBody() ([]byte, error)
	//encoder()[]byte // 最终数据
}

type IBodyDataWriter interface {
	//设置form数据并将content-type设置 为 application/x-www-form-urlencoded 或 multipart/form-data
	SetForm(values url.Values) error
	SetToForm(key, value string) error
	AddForm(key, value string) error
	// 会替换掉对应掉file信息，并且将content-type 设置为 multipart/form-data
	AddFile(key string, file *FileHeader) error
	//设置 multipartForm 数据并将content-type设置 为 multipart/form-data

	// 重置body，会清除掉未处理掉 form和file
	SetRaw(contentType string, body []byte)
}

type IBody interface {
	IBodyDataReader
	IBodyDataWriter
}

type IStatusGet interface {
	StatusCode() int
	Status() string
}

type IStatusSet interface {
	SetStatus(code int, status string)
}

type IRequestData interface {
	IBodyDataReader
	Method() string
	URL() *url.URL
	RequestURI() string
	Host() string
	RemoteAddr() string
	Scheme() string
}

// 原始请求数据的读
type IRequestReader interface {
	//ICookieReader
	IHeaderReader
	IRequestData
}

// 用于组装转发的request
type IRequest interface {
	//ICookieReader
	IHeaderReader
	IHeaderWriter
	ICookieWriter
	IBody
	Querys() url.Values
	TargetServer() string
	TargetURL() string
	Url() *url.URL
}

// 读取转发结果的response
type IResponseReader interface {
	//ICookieReader
	IHeaderReader
	IBodyGet
	IStatusGet
}

//// 单存储
//type IStore interface {
//	Set(value interface{})
//	Get() (value interface{})
//}

// 带优先的header
type IPriorityHeader interface {
	IHeaderReader // 读已经设置的header
	IHeaderWriter // 设置header
	// 非Priority的header会被 proxy 的同名项替换掉，
	Set() IHeader    // 这里设置的header会替换掉proxy的内容
	Append() IHeader // 这里设置的header会追加到proxy的内容
}

// 返回给client端的
type ResponseWriter interface {
	IPriorityHeader
	//ICookieReader // 已经设置的cookie
	ICookieWriter // 设置返回的cookie
	IStatusGet
	IStatusSet // 设置返回状态
	IBodySet   // 设置返回内容
	IBodyGet
}
