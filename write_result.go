package eosc

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/eolinker/eosc/log"
)

//PageInfo 页码信息
type PageInfo struct {
	ItemNum  int `json:"itemNum,"`
	Page     int `json:"page,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
	TotalNum int `json:"totalNum,"`
}

//SetPage 设置页码
func (p *PageInfo) SetPage(page, size, total int) *PageInfo {
	p.Page = page
	p.PageSize = size
	p.TotalNum = total
	return p
}

//NewItemNum 创建新的item
func NewItemNum(num int) *PageInfo {
	return &PageInfo{
		ItemNum:  num,
		TotalNum: num,
	}
}

//WriteResult 写响应
func WriteResult(w http.ResponseWriter, code string, resultType, resultKey string, result interface{}, pageInfo *PageInfo) {
	ret := map[string]interface{}{
		"statusCode": code,
	}
	if resultType != "" {
		ret["type"] = resultType
	}
	if result != nil {
		ret[resultKey] = result
	}
	if pageInfo != nil {
		ret["page"] = pageInfo
	}

	data, err := json.Marshal(ret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithFields(ret).Log(log.DebugLevel, err)
		return
	}

	i, err := w.Write(data)
	if err != nil {
		log.WithFields(ret).Log(log.DebugLevel, "write error:", err)
	} else {
		log.WithFields(ret).Log(log.DebugLevel, "write :", i)
	}
}

//WriteSuccessInData 返回包含于data中的成功信息
func WriteSuccessInData(w http.ResponseWriter, resultType, resultKey string, result interface{}) {
	WriteResultInfo(w, resultType, "IData", map[string]interface{}{resultKey: result})
}

//WriteSuccessNoData 返回无数据的成功信息
func WriteSuccessNoData(w http.ResponseWriter, resultType string) {
	WriteResultInfo(w, resultType, "", nil)
}

//WriteResultInfoWithPage 返回带页码信息的数据
func WriteResultInfoWithPage(w http.ResponseWriter, resultType, resultKey string, result interface{}, pageInfo *PageInfo) {
	WriteResult(w, "000000", resultType, resultKey, result, pageInfo)
}

//WriteResultInfo 返回信息
func WriteResultInfo(w http.ResponseWriter, resultType string, resultKey string, result interface{}) {

	if result != nil {
		if t := reflect.TypeOf(result); t.Kind() == reflect.Slice {

			WriteResultInfoWithPage(w, resultType, resultKey, result, NewItemNum(reflect.ValueOf(result).Len()))
			return
		}
	}
	WriteResultInfoWithPage(w, resultType, resultKey, result, nil)
}

//WriteResultInfoWithCode 返回带状态信息的响应
func WriteResultInfoWithCode(w http.ResponseWriter, code string, resultType, resultKey string, result interface{}) {
	WriteResult(w, code, resultType, resultKey, result, nil)
}

//WriteError 返回错误
func WriteError(w http.ResponseWriter, resultErr error) {

	w.Write([]byte(resultErr.Error()))
}
