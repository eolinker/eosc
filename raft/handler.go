package raft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-basic/uuid"

	"github.com/eolinker/eosc/log"
)

func (rc *Node) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if rc.transportHandler != nil {
		rc.transportHandler.ServeHTTP(writer, request)
		return
	}
	http.NotFound(writer, request)

}

// genHandler http请求处理
func (rc *Node) genHandler() http.Handler {
	sm := http.NewServeMux()
	// 其他节点加入集群的处理
	sm.HandleFunc("/raft/node/join", rc.joinHandler)

	sm.HandleFunc("/raft/node/info", rc.getNodeInfo)
	// 其他节点转发到leader的处理
	sm.HandleFunc("/raft/node/propose", rc.proposeHandler)

	sm.Handle("/", rc.transport.Handler())
	return sm
}

//getNodeInfo 获取集群信息
func (rc *Node) getNodeInfo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	joinData, err := decodeJoinRequest(body)
	if err != nil {
		writeError(w, "110001", "fail to parse join data", err.Error())
		return
	}
	if !rc.join {
		err = rc.UpdateHostInfo(joinData.Target)
		if err != nil {
			writeError(w, "110002", "fail to update host Info", err.Error())
			return
		}
	}
	writeSuccessResult(w, "", &JoinResponse{
		NodeSecret: &NodeSecret{
			ID:  rc.peers.Index() + 1,
			Key: uuid.New(),
		},
		Peer: rc.peers.GetAllPeers(),
	})
	return
}

// joinHandler 收到其他节点加入集群的处理
// 1、如果已经是集群模式，直接返回相关id，peer等信息方便处理
// 2、如果不是集群模式，先切换集群rc.changeCluster,再返回相关信息
// 3、该处理也可应用于集群节点crash后的重启
func (rc *Node) joinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	joinData, err := decodeJoinRequest(body)
	if err != nil {
		writeError(w, "110001", "fail to parse join data", err.Error())
		return
	}

	addr := fmt.Sprintf("%s://%s", joinData.Protocol, joinData.BroadcastIP)
	if joinData.BroadcastPort > 0 {
		addr = fmt.Sprintf("%s:%d", addr, joinData.BroadcastPort)
	}
	log.Infof("address %s apply join the cluster", addr)
	if id, exist := rc.peers.CheckExist(addr); exist {
		// 当前地址已经存在
		if id != joinData.NodeID {
			// ID错误
			writeError(w, "110004", "id and address do not match", "id and address do not match")
			return
		}
		writeSuccessResult(w, "", nil)
		return
	}

	node := &NodeInfo{
		NodeSecret: &NodeSecret{
			ID:  joinData.NodeID,
			Key: joinData.NodeKey,
		},
		BroadcastIP:   joinData.BroadcastIP,
		BroadcastPort: joinData.BroadcastPort,
		Addr:          addr,
		Protocol:      joinData.Protocol,
	}
	data, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	// 发送新增节点请求
	err = rc.AddNode(node.ID, data)
	if err != nil {
		writeError(w, "110003", "fail to add config error", err.Error())
		return
	}
	writeSuccessResult(w, "", nil)
	return
}

// proposeHandler 其他节点转发到leader的propose处理，由rc.Send触发
func (rc *Node) proposeHandler(w http.ResponseWriter, r *http.Request) {
	// 只有leader才会收到该消息

	defer r.Body.Close()

	isLeader, err := rc.isLeader()
	if err != nil {
		return
	}
	if !isLeader {
		writeError(w, "120001", "can not find leader", "can not find leader")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, "120002", "fail to read body", err.Error())
		return
	}

	msg, err := decodeProposeMsg(b)
	if err != nil {
		w.WriteHeader(503)
		writeError(w, "120003", "fail to parse propose message", err.Error())
		return
	}
	obj, data, err := rc.service.ProcessDataHandler(msg.Body)
	if err != nil {
		w.WriteHeader(503)
		writeError(w, "120004", "fail to send propose message", err.Error())
		return
	}
	err = rc.ProcessData(data)
	if err != nil {
		w.WriteHeader(503)
		writeError(w, "120005", "fail to send propose message", err.Error())
		return
	}
	err = rc.service.CommitHandler(data)
	if err != nil {
		w.WriteHeader(503)
		writeError(w, "120005", "fail to commit message", err.Error())
		return
	}
	writeTo(w, obj)
}

func encodeProposeMsg(from uint64, data []byte) ([]byte, error) {
	msg := &ProposeMsg{
		Body: data,
		From: from,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil

}
func decodeProposeMsg(data []byte) (*ProposeMsg, error) {
	msg := &ProposeMsg{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeResponse(data []byte) (*Response, error) {
	res := &Response{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func decodeJoinRequest(data []byte) (*JoinRequest, error) {
	joinData := new(JoinRequest)
	err := json.Unmarshal(data, joinData)
	if err != nil {
		return nil, err
	}
	return joinData, nil
}

func decodeJoinResponse(data []byte) (*JoinResponse, error) {
	res := new(JoinResponse)
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
