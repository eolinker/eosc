package raft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-basic/uuid"

	"github.com/eolinker/eosc/log"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/pkg/wait"

	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var defaultSnapshotCount uint64 = 10000
var snapshotCatchUpEntriesN uint64 = 10000

var retryFrequency time.Duration = 2000

// JoinCluster 新建一个加入已知集群的请求
// join和isCluster一定为true
// 1、应用于新建一个想要加入已知集群的节点，会向已知节点发送请求获取id等新建节点信息
// 已知节点如果还处于非集群模式，会先切换成集群模式
// 2、也可以用于节点crash后的重启处理
func JoinCluster(broadCastIP string, broadPort int, target string, service IService, count int) (*Node, error) {
	if count > 2 {
		return nil, errors.New("join error")
	}
	msg := JoinRequest{
		BroadcastIP:   broadCastIP,
		BroadcastPort: broadPort,
		Protocol:      "http",
		Target:        target,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	log.Info("dadqer")
	// 向集群中的某个节点发送要加入的请求
	resp, err := http.Post(fmt.Sprintf("%s/raft/node/join", target), "application/json;charset=utf-8", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	log.Info("asdqr")
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	err = json.Unmarshal(content, res)
	if err != nil {
		return nil, err
	}
	if res.Code == "000000" {
		resMsg := &JoinResponse{}
		log.Info("dqwqeer")
		data, _ := json.Marshal(res.Data)
		err = json.Unmarshal(data, resMsg)
		if err != nil {
			return nil, err
		}

		if resMsg.ResponseType != "join" {
			node := &NodeInfo{
				NodeSecret:    resMsg.NodeSecret,
				BroadcastIP:   broadCastIP,
				BroadcastPort: broadPort,
			}
			resMsg.Peer[node.ID] = node
			n := joinAndCreateRaft(node, service, resMsg.Peer)

			_, err = JoinCluster(broadCastIP, broadPort, target, service, count+1)
			if err != nil {
				return nil, err
			}
			return n, nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf(res.Err)

}

// joinAndCreateRaft 收到id，peer等信息后，新建并加入集群，新建日志文件等处理
func joinAndCreateRaft(node *NodeInfo, service IService, peers map[uint64]*NodeInfo) *Node {
	rc := NewNode(service)
	rc.nodeID = node.ID
	rc.waldir = fmt.Sprintf("eosc-%d", node.ID)
	rc.snapdir = fmt.Sprintf("eosc-%d-snap", node.ID)
	rc.join, rc.isCluster = true, true
	rc.nodeKey = uuid.New()
	rc.transport.ID = types.ID(rc.nodeID)
	rc.transport.Raft = rc
	rc.transport.LeaderStats = stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID)))

	for _, p := range peers {
		rc.peers.SetPeer(p.ID, p)
	}

	//rc.updateTransport <- true
	go rc.startRaft()
	return rc
}

//NewNode 新建raft节点
func NewNode(service IService) *Node {
	logger := zap.NewExample()
	rc := &Node{
		peers:           NewPeers(),
		Service:         service,
		snapCount:       defaultSnapshotCount,
		stopc:           make(chan struct{}),
		httpstopc:       make(chan struct{}),
		httpdonec:       make(chan struct{}),
		logger:          zap.NewExample(),
		waiter:          wait.New(),
		lead:            0,
		active:          true,
		updateTransport: make(chan bool, 1),
		transport: &rafthttp.Transport{
			Logger:             logger,
			ClusterID:          0x1000,
			ServerStats:        stats.NewServerStats("", ""),
			ErrorC:             make(chan error),
			DialRetryFrequency: rate.Every(2000 * time.Millisecond),
		},
	}

	rc.transport.Raft = rc

	return rc
}
