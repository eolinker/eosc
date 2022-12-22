package http_context

import (
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"time"

	"github.com/fasthttp/websocket"

	"github.com/eolinker/eosc/eocontext"
)

type IWebsocketContext interface {
	IHttpContext
	Upgrade() error
	SetUpstreamConn(conn *websocket.Conn)
	IsWebsocket() bool
}

type IHttpContext interface {
	eocontext.EoContext
	Request() IRequestReader // 读取原始请求
	Proxy() IRequest         // 读写转发请求
	Response() IResponse     // 处理返回结果，可读可写
	SendTo(address string, timeout time.Duration) error
	Proxies() []IProxy
	FastFinish()
}

type IHeaderReader interface {
	RawHeader() string
	GetHeader(name string) string
	Headers() http.Header
	Host() string
	GetCookie(key string) string
}

type IHeaderWriter interface {
	IHeaderReader
	SetHeader(key, value string)
	AddHeader(key, value string)
	DelHeader(key string)
	SetHost(host string)
}

type IResponseHeader interface {
	GetHeader(name string) string
	Headers() http.Header
	HeadersString() string
	SetHeader(key, value string)
	AddHeader(key, value string)
	DelHeader(key string)
}
type IBodyGet interface {
	GetBody() []byte
	BodyLen() int
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
	Files() (map[string][]*multipart.FileHeader, error)
	GetForm(key string) string
	GetFile(key string) (file []*multipart.FileHeader, has bool)
	RawBody() ([]byte, error)
}

type IBodyDataWriter interface {
	IBodyDataReader
	//设置form数据并将content-type设置 为 application/x-www-form-urlencoded 或 multipart/form-data
	SetForm(values url.Values) error
	SetToForm(key, value string) error
	AddForm(key, value string) error
	// 会替换掉对应掉file信息，并且将content-type 设置为 multipart/form-data
	AddFile(key string, file *multipart.FileHeader) error
	//设置 multipartForm 数据并将content-type设置 为 multipart/form-data
	// 重置body，会清除掉未处理掉 form和file
	SetRaw(contentType string, body []byte)
}

type IStatusGet interface {
	StatusCode() int
	ProxyStatusCode() int
	ProxyStatus() string
	Status() string
}

type IStatusSet interface {
	SetStatus(code int, status string)
	SetProxyStatus(code int, status string)
}

type IQueryReader interface {
	GetQuery(key string) string
	RawQuery() string
}

type IQueryWriter interface {
	IQueryReader
	SetQuery(key, value string)
	AddQuery(key, value string)
	DelQuery(key string)
	SetRawQuery(raw string)
}

type IURIReader interface {
	RequestURI() string
	Scheme() string
	RawURL() string
	Host() string
	Path() string
	IQueryReader
}

type IURIWriter interface {
	IURIReader
	IQueryWriter
	//SetRequestURI(uri string)
	SetPath(string)
	SetScheme(scheme string)
	SetHost(host string)
}

// 原始请求数据的读
type IRequestReader interface {
	Header() IHeaderReader
	Body() IBodyDataReader
	RemoteAddr() string
	RemotePort() string
	ReadIP() string
	ForwardIP() string
	URI() IURIReader
	Method() string
	String() string
	ContentLength() int
	ContentType() string
}

// 用于组装转发的request
type IRequest interface {
	Method() string
	ContentLength() int
	ContentType() string
	Header() IHeaderWriter
	Body() IBodyDataWriter
	URI() IURIWriter
	SetMethod(method string)
}

// IProxy 记录转发相关信息
type IProxy interface {
	IRequest
	StatusCode() int
	Status() string
	ProxyTime() time.Time
	ResponseLength() int
	ResponseTime() int64
}

// 返回给client端的
type IResponse interface {
	ResponseError() error
	ClearError()
	String() string
	IStatusGet
	IResponseHeader
	IStatusSet // 设置返回状态
	IBodySet   // 设置返回内容
	IBodyGet
	SetResponseTime(duration time.Duration)
	ResponseTime() time.Duration
	ContentLength() int
	ContentType() string
}
