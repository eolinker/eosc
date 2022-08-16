package open_api

import (
	"encoding/json"
	"github.com/eolinker/eosc/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OpenApiHandler interface {
	Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*EventResponse, body interface{})
}
type OpenApiHandleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*EventResponse, body interface{})

func (f OpenApiHandleFunc) Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, event []*EventResponse, body interface{}) {
	return f(req, params)
}
func CreateHandleFunc(handleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*EventResponse, body interface{})) httprouter.Handle {
	return CreateHandler(OpenApiHandleFunc(handleFunc))
}

func CreateHandler(handler OpenApiHandler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
		status, header, event, body := handler.Serve(req, param)

		response := &Response{
			StatusCode: status,
			Data:       nil,
			Event:      event,
		}
		if header == nil {
			header = make(http.Header)
			header.Set("content-type", "application/json")
		} else {
			if header.Get("content-type") == "" {
				header.Set("content-type", "application/json")
			}
		}
		response.Header = header
		if body != nil {
			switch d := body.(type) {
			case error:
				response.Data = []byte(d.Error())
				log.Debug("handler write err:", d)

			case string:
				response.Data = []byte(d)
				log.Debug("handler write string:", string(d))

			case []byte:
				response.Data = d
				log.Debug("handler write []byte:", string(d))
			default:
				response.Data, _ = json.Marshal(body)
				log.Debug("handler write default:", string(response.Data))
			}
		}
		data, _ := json.Marshal(response)

		w.WriteHeader(200)
		w.Write(data)
	}
}
