package open_api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OpenApiHandler interface {
	Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, event *EventResponse, body interface{})
}
type OpenApiHandleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, event *EventResponse, body interface{})

func (f OpenApiHandleFunc) Serve(req *http.Request, params httprouter.Params) (status int, header http.Header, event *EventResponse, body interface{}) {
	return f(req, params)
}
func CreateHandleFunc(handleFunc func(req *http.Request, params httprouter.Params) (status int, header http.Header, event *EventResponse, body interface{})) httprouter.Handle {
	return CreateHandler(OpenApiHandleFunc(handleFunc))
}

func CreateHandler(handler OpenApiHandler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
		status, header, event, body := handler.Serve(req, param)

		response := &Response{
			StatusCode: status,
			Header:     header,
			Data:       nil,
			Event:      event,
		}

		if body != nil {
			switch d := body.(type) {
			case error:
				response.Data = []byte(d.Error())
			case string:
				response.Data = []byte(d)
			case []byte:
				response.Data = d
			default:
				response.Data, _ = json.Marshal(body)
			}
		}
		data, _ := json.Marshal(response)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
	}
}
