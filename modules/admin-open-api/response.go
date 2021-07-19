package admin_open_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResultError(w http.ResponseWriter,status int, err error)  {
	w.WriteHeader(500)
	fmt.Fprintf(w,"%s",err.Error())
}

func writeResult(w http.ResponseWriter,v interface{})  {
	data,err:=json.Marshal(v)
	if err!= nil{
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}