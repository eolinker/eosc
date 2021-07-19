package eosc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponse(w http.ResponseWriter,data interface{})  {
	if body,ok:= data.([]byte);ok{
		w.WriteHeader(200)
		w.Write(body)
		return
	}
	if body,err:=json.Marshal(data); err!= nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,"Internal Server Error:%s",err.Error())
		return
	}else{
		w.WriteHeader(200)
		w.Write(body)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string)  {
	w.WriteHeader(statusCode)
	fmt.Fprint(w,message)
}