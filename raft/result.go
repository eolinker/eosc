package raft

import (
	"encoding/json"
	"net/http"
)

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
	data, _ := json.Marshal(result)
	w.Write(data)
}
