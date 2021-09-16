package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/eolinker/eosc/log"
)

//getNodeInfoRequest 发送获取节点信息请求
func getNodeInfoRequest(address string, data []byte) (*JoinResponse, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/info", address))
	if err != nil {
		log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
		return nil, err
	}
	// 向集群中的某个节点发送要加入的请求
	resp, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	res, err := decodeResponse(content)
	if err != nil {
		return nil, err
	}
	if res.Code == "000000" {

		data, _ := json.Marshal(res.Data)
		resMsg, err := decodeJoinResponse(data)
		if err != nil {
			return nil, err
		}
		return resMsg, nil
	}
	return nil, fmt.Errorf(res.Err)
}

//joinClusterRequest 发送加入集群请求
func joinClusterRequest(address string, data []byte) error {
	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join", address))
	if err != nil {
		log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
		return err
	}
	// 向集群中的某个节点发送要加入的请求
	resp, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res, err := decodeResponse(content)
	if err != nil {
		return err
	}
	if res.Code == "000000" {
		return nil
	}
	return fmt.Errorf(res.Err)
}
