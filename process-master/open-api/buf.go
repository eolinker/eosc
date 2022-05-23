package open_api

import (
	"bytes"
	"net/http"
)

type _ProxyWriterBuffer struct {
	buf        bytes.Buffer
	statusCode int
	header     http.Header
}

func NewTemplateWriter() *_ProxyWriterBuffer {
	return &_ProxyWriterBuffer{
		statusCode: 200,
		header:     make(http.Header),
	}
}
func (t *_ProxyWriterBuffer) WriteTo(w http.ResponseWriter) {

	t.WriteHeaderTo(w)
	w.WriteHeader(t.statusCode)

	t.buf.WriteTo(w)
}
func (t *_ProxyWriterBuffer) WriteHeaderTo(w http.ResponseWriter) {
	for k := range t.header {
		w.Header().Set(k, t.header.Get(k))
	}
}
func (t *_ProxyWriterBuffer) Header() http.Header {
	return t.header
}

func (t *_ProxyWriterBuffer) Write(bytes []byte) (int, error) {
	return t.buf.Write(bytes)
}

func (t *_ProxyWriterBuffer) WriteHeader(statusCode int) {
	t.statusCode = statusCode
}
