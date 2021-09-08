package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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
func JoinCluster(broadCastIP string, broadPort int, target string, service IService) (*Node, error) {
	msg := JoinRequest{
		BroadcastIP:   broadCastIP,
		BroadcastPort: broadPort,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// 向集群中的某个节点发送要加入的请求
	resp, err := http.Post(fmt.Sprintf("%s/raft/join", target), "application/json;charset=utf-8", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
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
		data, _ := json.Marshal(res.Data)
		err = json.Unmarshal(data, resMsg)
		if err != nil {
			return nil, err
		}
		log.Infof("receive join message id(%d), peers number(%d)", resMsg.ID, len(resMsg.Peer))
		node := &NodeInfo{
			NodeSecret:    resMsg.NodeSecret,
			BroadcastIP:   broadCastIP,
			BroadcastPort: broadPort,
		}
		resMsg.Peer[node.ID] = node
		return joinAndCreateRaft(node, service, resMsg.Peer), nil
	} else {
		return nil, fmt.Errorf(res.Msg)
	}
}

// joinAndCreateRaft 收到id，peer等信息后，新建并加入集群，新建日志文件等处理
func joinAndCreateRaft(node *NodeInfo, service IService, peers map[uint64]*NodeInfo) *Node {
	rc := &Node{
		nodeID:    node.ID,
		Service:   service,
		join:      true,
		isCluster: true,
		// 日志文件目前直接以id命名初始化
		waldir:    fmt.Sprintf("eosc-%d", node.ID),
		snapdir:   fmt.Sprintf("eosc-%d-snap", node.ID),
		snapCount: defaultSnapshotCount,
		stopc:     make(chan struct{}),
		httpstopc: make(chan struct{}),
		httpdonec: make(chan struct{}),
		logger:    zap.NewExample(),
		waiter:    wait.New(),
		lead:      0,
		active:    true,
	}

	// 确保peer一定有节点本身
	rc.peers = NewPeers()
	for _, p := range peers {
		rc.peers.SetPeer(p.ID, p)
	}
	// 创建并启动 transport 实例，该负责节点之间的网络通信，
	// 非集群模式下主要是为了listener的Handler处理，监听join请求，此时transport尚未start
	rc.transport = &rafthttp.Transport{
		Logger:             rc.logger,
		ID:                 types.ID(rc.nodeID),
		ClusterID:          0x1000,
		Raft:               rc,
		ServerStats:        stats.NewServerStats("", ""),
		LeaderStats:        stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID))),
		ErrorC:             make(chan error),
		DialRetryFrequency: rate.Every(retryFrequency * time.Millisecond),
	}

	// raft启动
	//go rc.serveRaft()
	go rc.startRaft()
	return rc
}

//// CreateRaftNode 初始化节点
//// peers至少会包括节点本身，如果join为true，isCluster也默认为true
//// 已经切换到集群模式下的节点不会回到非集群模式，除非改变节点ID或删除相关日志文件
//// 1、创建非集群节点,isCluster为false，此时peers可为空
//// 2、创建集群节点，isCluster为true，若此时peers为空(或仅有节点本身)，表示该集群仅有一个节点
//// peers也可以是其余集群中的其他节点，表示这是一个多节点集群，此时其他节点也需通过同样的配置和方式启动，
//// 推荐使用JoinCluster来新建多节点集群节点
//// 3、创建加入已知集群的节点，join为true，isCluster为true，此时peers需包括其他节点地址，推荐使用JoinCluster来新建非单点集群节点
//func CreateRaftNode(id int, host string, service IService, peers string, keys string, join bool, isCluster bool) (*Node, error) {
//	rc := &Node{
//		nodeID:    uint64(id),
//		Service:   service,
//		join:      join,
//		waldir:    fmt.Sprintf("eosc-%d", id), // 日志文件路径
//		snapdir:   fmt.Sprintf("eosc-%d-snap", id),
//		snapCount: defaultSnapshotCount,
//		stopc:     make(chan struct{}),
//		httpstopc: make(chan struct{}),
//		httpdonec: make(chan struct{}),
//		logger:    zap.NewExample(),
//		waiter:    wait.New(),
//		lead:      0,
//		active:    true,
//	}
//	// 建过集群的节点不能再换回去（暂时采用该方案）
//	if rc.join || wal.Exist(rc.waldir) {
//		rc.isCluster = true
//	}
//
//	log.Infof("current mode is cluster %v.", rc.isCluster)
//
//	peerList, err := Adjust(rc.nodeID, host, peers, keys, rc.isCluster)
//	if err != nil {
//		log.Error(err.Error())
//		return nil, err
//	}
//	rc.peers = NewPeers()
//
//	// 创建并启动 transport 实例，负责节点之间的网络通信，
//	// 非集群模式下主要是为了listener的Handler处理，监听join请求，此时transport尚未start
//	rc.transport = &rafthttp.Transport{
//		Logger:             rc.logger,
//		ID:                 types.ID(rc.nodeID),
//		ClusterID:          0x1000,
//		Raft:               rc,
//		ServerStats:        stats.NewServerStats("", ""),
//		LeaderStats:        stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID))),
//		ErrorC:             make(chan error),
//		DialRetryFrequency: rate.Every(2000 * time.Millisecond),
//	}
//
//	// 监听节点端口，用transport处理节点通信，此时这种情况下只是监听join
//	go rc.serveRaft()
//
//	// 集群模式下启动节点，已经是认为是集群的节点（有日志文件存在）也会启动集群模式
//	if isCluster || wal.Exist(rc.waldir) {
//		go rc.startRaft()
//	}
//	return rc, nil
//}

func NewNode(service IService) (*Node, error) {
	rc := &Node{
		Service:   service,
		snapCount: defaultSnapshotCount,
		stopc:     make(chan struct{}),
		httpstopc: make(chan struct{}),
		httpdonec: make(chan struct{}),
		logger:    zap.NewExample(),
		waiter:    wait.New(),
		lead:      0,
		active:    true,
	}

	// 创建并启动 transport 实例，负责节点之间的网络通信，
	// 非集群模式下主要是为了listener的Handler处理，监听join请求，此时transport尚未start
	rc.transport = &rafthttp.Transport{
		Logger:             rc.logger,
		ID:                 types.ID(rc.nodeID),
		ClusterID:          0x1000,
		Raft:               rc,
		ServerStats:        stats.NewServerStats("", ""),
		LeaderStats:        stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID))),
		ErrorC:             make(chan error),
		DialRetryFrequency: rate.Every(2000 * time.Millisecond),
	}
	rc.peers = NewPeers()
	// 监听节点端口，用transport处理节点通信，此时这种情况下只是监听join
	//go rc.serveRaft()

	return rc, nil
}

//// Adjust 参数校验与调整
//func Adjust(id uint64, host string, peers string, keys string, isCluster bool) (map[uint64]string, error) {
//	peerList := make(map[uint64]string)
//	peerList[id] = host
//	// 非集群模式不需要peer列表，此时peer列表仅有节点本身
//	if !isCluster {
//		return peerList, nil
//	}
//	clusterList := strings.Split(peers, ",")
//	keyList := strings.Split(keys, ",")
//	if len(keyList) != len(clusterList) {
//		return nil, fmt.Errorf("the length of keys is %d while the length of peers is %d", len(keyList), len(clusterList))
//	} else {
//		for i, key := range keyList {
//			k, err := strconv.Atoi(key)
//			if err != nil {
//				return nil, err
//			}
//			if _, ok := peerList[uint64(k)]; !ok {
//				peerList[uint64(k)] = clusterList[i]
//			}
//		}
//	}
//	return peerList, nil
//}
