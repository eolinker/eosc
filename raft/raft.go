package raft

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.etcd.io/etcd/raft/v3/raftpb"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/pkg/wait"
	"go.etcd.io/etcd/raft/v3"

	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var defaultSnapshotCount uint64 = 10000
var snapshotCatchUpEntriesN uint64 = 10000

var retryFrequency time.Duration = 2000

// JoinCluster 新建一个加入已知集群的请求
// 1、应用于新建一个想要加入已知集群的节点，会向已知节点发送请求获取id等新建节点信息
// 已知节点如果还处于非集群模式，会先切换成集群模式
// 2、也可以用于节点crash后的重启处理
func JoinCluster(rc *Node, broadCastIP string, broadPort int, address string, protocol string) error {
	// 判断是否已经在一个多节点集群中
	if rc.peers.GetPeerNum() > 1 || rc.join {
		return fmt.Errorf("This node has joined the cluster")
	}
	msg := JoinRequest{
		BroadcastIP:   broadCastIP,
		BroadcastPort: broadPort,
		Protocol:      protocol,
		Target:        address,
	}
	data, _ := json.Marshal(msg)

	resp, err := getNodeInfoRequest(address, data)
	if err != nil {
		return err
	}

	rc.nodeID = resp.ID
	rc.nodeKey = resp.Key
	rc.broadcastPort = broadPort
	rc.broadcastIP = broadCastIP
	rc.protocol = protocol

	err = rc.joinInit()
	if err != nil {
		return err
	}
	err = startRaft(rc, resp.Peer)
	if err != nil {
		return err
	}

	msg.NodeID = resp.ID
	msg.NodeKey = resp.Key
	data, _ = json.Marshal(msg)
	err = joinClusterRequest(address, data)
	if err != nil {
		return err
	}
	return nil
}

// JoinInit 加入集群前的初始化
func (rc *Node) joinInit() error {
	// 关闭当前单例集群服务
	rc.stop()
	// 删除旧的日志文件
	err := rc.removeWalFile()
	if err != nil {
		return err
	}
	rc.join = true
	return nil
}

// startRaft 收到id，peer等信息后，新建并加入集群，新建日志文件等处理
func startRaft(rc *Node, peers map[uint64]*NodeInfo) error {
	rc.waldir = fmt.Sprintf("eosc-%d", rc.nodeID)
	rc.snapdir = fmt.Sprintf("eosc-%d-snap", rc.nodeID)
	rc.transport.ID = types.ID(rc.nodeID)
	rc.transport.Raft = rc
	rc.transport.LeaderStats = stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID)))
	rc.transportHandler = rc.genHandler()
	rc.stopc = make(chan struct{})

	rc.peers.SetPeer(rc.nodeID, &NodeInfo{
		NodeSecret: &NodeSecret{
			ID:  rc.nodeID,
			Key: rc.nodeKey,
		},
		BroadcastIP:   rc.broadcastIP,
		BroadcastPort: rc.broadcastPort,
		Protocol:      rc.protocol,
	})
	for _, p := range peers {
		rc.peers.SetPeer(p.ID, p)
	}
	rc.writeConfig()
	return rc.startRaft()
}

//NewNode 新建raft节点
func NewNode(service IService) (*Node, error) {
	//fileName := fmt.Sprintf("%s_node.args", eosc_args.AppName())
	//// 判断是否存在nodeID，若存在，则当作旧节点处理，加入集群
	//cfg := eosc_args.NewConfig(fileName)
	//cfg.ReadFile(fileName)
	// 均已node_id为1启动,作为单例集群
	//nodeID, _ := strconv.Atoi(cfg.GetDefault(eosc_args.NodeID, "1"))
	//nodeKey := cfg.GetDefault(eosc_args.NodeKey, "")
	logger, _ := zap.NewProduction()
	rc := &Node{
		peers:     NewPeers(),
		service:   service,
		snapCount: defaultSnapshotCount,
		stopc:     make(chan struct{}),
		logger:    logger,
		waiter:    wait.New(),
		lead:      0,
		transport: &rafthttp.Transport{
			Logger:             logger,
			ClusterID:          0x1000,
			ServerStats:        stats.NewServerStats("", ""),
			ErrorC:             make(chan error),
			DialRetryFrequency: rate.Every(2000 * time.Millisecond),
		},
	}
	rc.readConfig()

	err := startRaft(rc, nil)
	if err != nil {
		return nil, err
	}
	service.SetRaft(rc)
	return rc, nil
}

func (rc *Node) ProcessData(data []byte) error {

	m := &Message{
		From: rc.nodeID,
		Type: PROPOSE,
		Data: data,
	}
	b, err := m.Encode()
	if err != nil {
		return err
	}
	return rc.node.Propose(context.TODO(), b)
}
func (rc *Node) Process(ctx context.Context, m raftpb.Message) error {
	if rc.node == nil {
		return nil
	}
	return rc.node.Step(ctx, m)
}
func (rc *Node) IsIDRemoved(id uint64) bool { return false }
func (rc *Node) ReportUnreachable(id uint64) {
	if rc.node == nil {
		return
	}
	rc.node.ReportUnreachable(id)
}
func (rc *Node) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	rc.node.ReportSnapshot(id, status)
}
