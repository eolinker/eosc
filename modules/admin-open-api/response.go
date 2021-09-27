package admin_open_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eolinker/eosc/log"
)

func writeResultError(w http.ResponseWriter, status int, err error) {
	msg := err.Error()
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", msg)

	log.Infof("write error to client:%s", msg)
}

func writeResult(w http.ResponseWriter, v interface{}) {
	if v == nil {
		w.Write([]byte("{}"))
	}
	data, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}
