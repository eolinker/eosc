package etcdRaft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (e *EtcdServer) sendJoinRequest(target string, addr []string) (string, map[string][]string, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join", target))
	if err != nil {
		log.Printf("fail to join: addr is %s, error is %s", target, err.Error())
		return "", nil, nil
	}
	msg := joinRequest{
		Addr: addr,
	}
	data, _ := json.Marshal(msg)
	// 向集群中的某个节点发送要加入的请求
	resp, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return "", nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	res, err := decodeResponse(content)
	if err != nil {
		return "", nil, err
	}
	if res.Code != "000000" {
		return "", nil, fmt.Errorf(res.Msg)
	}
	data, _ = json.Marshal(res.Data)
	result := new(joinResponse)
	err = json.Unmarshal(data, result)
	if err != nil {
		return "", nil, err
	}
	return result.Name, result.Members, nil
}

//getNodeInfoRequest 发送获取节点信息请求
//func getNodeInfoRequest(rc *Node, ip string, port int, protocol, address string) (*JoinResponse, error) {
//	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join/try", address))
//	if err != nil {
//		log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
//		return nil, err
//	}
//	queries := url.Values{}
//	queries.Set("broadcast_ip", ip)
//	queries.Set("broadcast_port", strconv.Itoa(port))
//	queries.Set("protocol", protocol)
//	queries.Set("address", address)
//	//data, _ := json.Marshal(msg)
//	uri.RawQuery = queries.Encode()
//	resp, err := http.Get(uri.String())
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	content, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	res, err := decodeResponse(content)
//	if err != nil {
//		return nil, err
//	}
//
//	if res.Code == "000000" {
//
//		data, _ := json.Marshal(res.Data)
//		res := new(SNResponse)
//		err := json.Unmarshal(data, res)
//		if err != nil {
//			return nil, err
//		}
//		rc.lastSN = res.SN
//	}
//
//	// 向集群中的某个节点发送要加入的请求
//
//	resp, err = http.Post(uri.String(), "application/json;charset=utf-8", strings.NewReader(""))
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	content, err = ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	res, err = decodeResponse(content)
//	if err != nil {
//		return nil, err
//	}
//	if res.Code == "000000" {
//
//		data, _ := json.Marshal(res.Data)
//		resMsg, err := decodeJoinResponse(data)
//		if err != nil {
//			return nil, err
//		}
//		return resMsg, nil
//	}
//	return nil, fmt.Errorf(res.Err)
//}
//
////joinClusterRequest 发送加入集群请求
//func joinClusterRequest(id uint64, key string, ip string, port int, protocol, address string) error {
//	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join", address))
//	if err != nil {
//		log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
//		return err
//	}
//	msg := JoinRequest{
//		NodeID:        id,
//		NodeKey:       key,
//		BroadcastIP:   ip,
//		BroadcastPort: port,
//		Protocol:      protocol,
//		Target:        address,
//	}
//	data, _ := json.Marshal(msg)
//	// 向集群中的某个节点发送要加入的请求
//	resp, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBuffer(data))
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	content, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//	res, err := decodeResponse(content)
//	if err != nil {
//		return err
//	}
//	if res.Code == "000000" {
//		return nil
//	}
//	return fmt.Errorf(res.Err)
//}
//
////callbackSNRequest
//func callbackSNRequest(address string) (string, error) {
//	uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join/callback", address))
//	if err != nil {
//		log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
//		return "", err
//	}
//
//	// 向集群中的某个节点发送要加入的请求
//	resp, err := http.Get(uri.String())
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	content, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	res, err := decodeResponse(content)
//	if err != nil {
//		return "", err
//	}
//	if res.Code == "000000" {
//		data, _ := json.Marshal(res.Data)
//		resMsg, err := decodeSNResponse(data)
//		if err != nil {
//			return "", err
//		}
//		return resMsg.SN, err
//	}
//	return "", fmt.Errorf(res.Err)
//}
