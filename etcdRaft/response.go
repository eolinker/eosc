package etcdRaft

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"error,omitempty"`
}
type SNRequest struct {
	BroadcastIP   string `json:"broadcast_ip"`
	BroadcastPort int    `json:"broadcast_port"`
	Protocol      string `json:"protocol"`
	Target        string `json:"target"`
}
type SNResponse struct {
	SN string `json:"lastSN"`
}

// writeSuccessResult 返回成功结果
func writeSuccessResult(w http.ResponseWriter, key string, value interface{}) {
	result := &Response{
		Code: "000000",
		Msg:  "success",
	}
	if value != nil {
		if key == "" {
			result.Data = value
		} else {
			result.Data = map[string]interface{}{
				key: value,
			}
		}
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

// writeError 返回失败结果
func writeError(w http.ResponseWriter, code string, msg, errInfo string) {
	result := &Response{
		Code: code,
		Msg:  msg,
		Err:  errInfo,
	}
	writeTo(w, result)
}

func writeTo(w http.ResponseWriter, obj interface{}) {
	if data, ok := obj.([]byte); ok {
		w.Write(data)
		return
	}
	data, _ := json.Marshal(obj)
	w.Write(data)

}
