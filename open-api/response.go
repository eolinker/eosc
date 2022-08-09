package open_api

import (
	"fmt"
	"github.com/eolinker/eosc"
	"net/http"
)

type Response struct {
	StatusCode int              `json:"status"`
	Header     http.Header      `json:"header"`
	Data       []byte           `json:"data"`
	Event      []*EventResponse `json:"Event"`
}

type EventResponse struct {
	Event     string `json:"event"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Data      []byte `json:"data"`
}

func (e *EventResponse) String() string {
	switch e.Event {
	case eosc.EventSet:
		return fmt.Sprint(e.Event, " ", e.Namespace, "[", e.Key, "]=", string(e.Data))
	case eosc.EventDel:
		return fmt.Sprint(e.Event, " ", e.Namespace, "[", e.Key, "]")
	default:
		return e.Event
	}
}

type OpenApiProxyResponse struct {
	RespData interface{}    `json:"resp_data"`
	Event    *EventResponse `json:"event"`
}
