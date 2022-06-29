package etcdRaft

import (
	"encoding/json"
	"net/http"
)

type joinRequest struct {
	Addr []string `json:"addr"`
}
type joinResponse struct {
	Members map[string][]string `json:"members"`
	Name    string              `json:"name"`
}

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
//type SNRequest struct {
//	BroadcastIP   string `json:"broadcast_ip"`
//	BroadcastPort int    `json:"broadcast_port"`
//	Protocol      string `json:"protocol"`
//	Target        string `json:"target"`
//}
//type SNResponse struct {
//	SN string `json:"lastSN"`
//}

// writeSuccessResult 返回成功结果
func writeSuccessResult(w http.ResponseWriter, value interface{}) {
	result := &Response{
		Code: "000000",
		Msg:  "success",
		Data: value,
	}
	writeTo(w, result)
}

// writeError 返回失败结果
func writeError(w http.ResponseWriter, code string, errInfo string) {
	result := &Response{
		Code: code,
		Msg:  errInfo,
	}
	writeTo(w, result)
}

func writeTo(w http.ResponseWriter, obj interface{}) {
	if data, ok := obj.([]byte); ok {
		_, _ = w.Write(data)
		return
	}
	data, _ := json.Marshal(obj)
	_, _ = w.Write(data)
}

func decodeResponse(data []byte) (*Response, error) {
	res := new(Response)
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
