package open_api

import (
	"encoding/json"
	"github.com/eolinker/eosc/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OpenApiHandler interface {
	Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{})
}
type OpenApiHandleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{})

func (f OpenApiHandleFunc) Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	return f(req, params)
}
func CreateHandleFunc(handleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{})) httprouter.Handle {
	return CreateHandler(OpenApiHandleFunc(handleFunc))
}

func CreateHandler(handler OpenApiHandler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
		status, header, body := handler.Serve(req, param)

		if header == nil {
			header = make(http.Header)
			header.Set("content-type", "application/json")
		} else {
			if header.Get("content-type") == "" {
				header.Set("content-type", "application/json")
			}
		}
		var data []byte
		if body != nil {
			switch d := body.(type) {
			case error:
				data = []byte(d.Error())
				log.Debug("handler write err:", d)

			case string:
				data = []byte(d)
				log.Debug("handler write string:", string(d))

			case []byte:
				data = d
				log.Debug("handler write []byte:", string(d))
			default:
				data, _ = json.Marshal(body)
				log.Debug("handler write default:", string(data))
			}
		}

		for k, vs := range header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(status)

		w.Write(data)
	}
}
