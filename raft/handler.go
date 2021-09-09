package raft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-basic/uuid"

	"github.com/eolinker/eosc/log"
)

// Handler http请求处理
func (rc *Node) Handler() http.Handler {
	sm := http.NewServeMux()
	// 其他节点加入集群的处理
	sm.HandleFunc("/raft/join", rc.joinHandler)
	// 其他节点转发到leader的处理
	sm.HandleFunc("/raft/propose", rc.proposeHandler)

	//sm.HandleFunc("/raft/status", rc.proposeHandler)
	sm.Handle("/raft", rc.transport.Handler())
	return sm
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
	var joinData JoinRequest
	err = json.Unmarshal(body, &joinData)
	if err != nil {
		writeError(w, "110001", "fail to parse join data", err.Error())
		return
	}

	// 先判断是不是集群模式
	// 是的话返回要加入的相关信息
	// 不是的话先切换集群模式，再初始化startRaft()，再返回加入的相关信息
	if !rc.isCluster {
		// 非集群模式，先本节点切换成集群模式
		err = rc.changeCluster(joinData.Target)
		if err != nil {
			writeError(w, "110002", "fail to change cluster", err.Error())
			return
		}
		writeSuccessResult(w, "", &JoinResponse{
			NodeSecret: &NodeSecret{
				ID:  rc.peers.Index() + 1,
				Key: uuid.New(),
			},
			Peer:         rc.peers.GetAllPeers(),
			ResponseType: "cluster",
		})
		return
	}
	addr := fmt.Sprintf("%s://%s", joinData.Protocol, joinData.BroadcastIP)
	if joinData.BroadcastPort > 0 {
		addr = fmt.Sprintf("%s:%d", addr, joinData.BroadcastPort)
	}
	log.Infof("address %s apply join the cluster", addr)
	// 切换完了，开始新增对应节点并返回新增条件信息
	if id, exist := rc.peers.CheckExist(addr); exist {
		info, _ := rc.peers.GetPeerByID(id)
		writeSuccessResult(w, "", &JoinResponse{
			NodeSecret: &NodeSecret{
				ID:  info.ID,
				Key: info.Key,
			},
			Peer:         rc.peers.GetAllPeers(),
			ResponseType: "join",
		})
		return
	}

	node := &NodeInfo{
		NodeSecret: &NodeSecret{
			ID:  rc.peers.Index() + 1,
			Key: uuid.New(),
		},
		BroadcastIP:   joinData.BroadcastIP,
		BroadcastPort: joinData.BroadcastPort,
		Addr:          addr,
		Protocol:      joinData.Protocol,
	}
	data, _ := json.Marshal(node)
	// 已经是集群了，发送新增节点的消息后返回
	err = rc.AddConfigChange(node.ID, data)
	if err != nil {
		writeError(w, "110003", "fail to add config error", err.Error())
		return
	}
	writeSuccessResult(w, "", &JoinResponse{
		NodeSecret: &NodeSecret{
			ID:  node.ID,
			Key: node.Key,
		},
		Peer:         rc.peers.GetAllPeers(),
		ResponseType: "join",
	})
	return
}

// proposeHandler 其他节点转发到leader的propose处理，由rc.Send触发
func (rc *Node) proposeHandler(w http.ResponseWriter, r *http.Request) {
	// 只有leader才会收到该消息
	_, isLeader, err := rc.getLeader()
	if err != nil {
		writeError(w, "120001", "can not find leader", err.Error())
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
	defer r.Body.Close()

	msg := &ProposeMsg{}
	err = json.Unmarshal(b, msg)
	if err != nil {
		writeError(w, "120003", "fail to parse propose message", err.Error())
		return
	}
	log.Infof("receive propose request from node(%d)", msg.From)
	err = rc.Send(msg.Cmd, msg.Data)
	if err != nil {
		writeError(w, "120004", "fail to send propose message", err.Error())
		return
	}
	writeSuccessResult(w, "", nil)
}
