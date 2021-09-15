package raft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	eosc_args "github.com/eolinker/eosc/eosc-args"

	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/pkg/wait"
	"go.etcd.io/etcd/raft/v3"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.etcd.io/etcd/server/v3/wal"
	"go.uber.org/zap"
)

// raft节点结构
type Node struct {
	// 节点ID
	nodeID uint64

	// eosc 服务相关
	service IService

	// 节点相关
	node raft.Node

	nodeKey string

	broadcastIP string

	broadcastPort int

	protocol string

	// 当前leader
	lead uint64

	// 节点状态
	confState raftpb.ConfState

	// 集群相关
	join bool
	// 集群模式
	isCluster bool
	// peers列表
	peers *Peers

	// 快照相关
	snapshotter *snap.Snapshotter
	// 达到snapCount时本地保存快照
	snapCount uint64
	// 快照记录的索引
	snapshotIndex uint64
	snapdir       string // path to snapshot directory

	// 日志与内存
	raftStorage  *MemoryStorage
	wal          *wal.WAL
	waldir       string // path to WAL directory
	appliedIndex uint64

	// 与其他节点通信
	transport *rafthttp.Transport
	stopc     chan struct{} // signals proposal channel closed
	//httpstopc chan struct{} // signals http server to shutdown
	//httpdonec chan struct{} // signals http server shutdown complete

	// 日志相关，后续改为eosc_log
	logger           *zap.Logger
	waiter           wait.Wait
	active           bool
	updateTransport  chan bool
	transportHandler http.Handler
}

func (rc *Node) NodeID() uint64 {
	return rc.nodeID
}

func (rc *Node) NodeKey() string {
	return rc.nodeKey
}

//writeConfig 将raft节点的运行配置写进文件中
func (rc *Node) writeConfig() {
	cfg := eosc_args.NewConfig(fmt.Sprintf("%s_node.args", eosc_args.AppName()))
	cfg.Set(eosc_args.IsCluster, "true")
	cfg.Set(eosc_args.NodeID, strconv.Itoa(int(rc.nodeID)))
	cfg.Set(eosc_args.NodeKey, strconv.Itoa(int(rc.nodeID)))
	cfg.Set(eosc_args.BroadcastIP, rc.broadcastIP)
	cfg.Set(eosc_args.Port, strconv.Itoa(rc.broadcastPort))
	cfg.Save()
}

func (rc *Node) clearConfig() {
	cfg := eosc_args.NewConfig(fmt.Sprintf("%s_node.args", eosc_args.AppName()))
	cfg.Save()
	rc.nodeID = 0
	rc.nodeKey = ""
	rc.broadcastIP = ""
	rc.broadcastPort = 0
	rc.active, rc.join, rc.isCluster = false, false, false
	rc.transportHandler = rc.genHandler()
}

// startRaft 启动raft服务，在集群模式下启动或join模式下启动
// 非集群模式下启动的节点不会调用该start函数
func (rc *Node) startRaft() {
	log.Info("start raft service")

	// 判断快照文件夹是否存在，不存在则创建
	if !fileutil.Exist(rc.snapdir) {
		if err := os.Mkdir(rc.snapdir, 0750); err != nil {
			log.Fatalf("eosc: cannot create dir for snapshot (%v)", err)
		}
	}
	// 新建快照管理
	rc.snapshotter = snap.New(zap.NewExample(), rc.snapdir)

	// 将日志wal重写入raftNode实例中，读取快照和日志，并初始化raftStorage
	rc.wal = rc.replayWAL()

	// 集群模式下启动节点的时候，重新reload快照到service中
	// TODO 非集群想要切换成集群的时候，要么这里做进一步校验，要么切换前先存好快照和日志
	err := rc.ReadSnap(rc.snapshotter)
	if err != nil {
		log.Info("reload snap to service error:", err)
	}

	// 节点配置
	c := &raft.Config{
		ID:                        rc.nodeID,
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   rc.raftStorage,
		MaxSizePerMsg:             1024 * 1024,
		MaxInflightMsgs:           256,
		MaxUncommittedEntriesSize: 1 << 30,
	}
	peersList := rc.peers.GetAllPeers()
	// 判断是否有日志文件目录
	oldWal := wal.Exist(rc.waldir)
	// 启动node节点
	if rc.join || oldWal {
		// 选择加入集群或已有日志消息（曾经切换过集群模式）
		rc.node = raft.RestartNode(c)
	} else {
		// 启动节点时添加peers
		peers := make([]raft.Peer, 0, rc.peers.GetPeerNum())
		for id := range peersList {
			peers = append(peers, raft.Peer{ID: id})
		}
		// 新开一个集群
		rc.node = raft.StartNode(c, peers)
	}
	// 开启节点间通信
	// 通信实例开始运行
	err = rc.transport.Start()
	if err != nil {
		log.Info("transport start error:", err)
	}
	// 与集群中的其他节点建立通信

	for k, v := range peersList {
		// transport加入peer列表，节点本身不添加
		if k != rc.nodeID {
			fmt.Println("v addr is", v.Addr)
			rc.transport.AddPeer(types.ID(k), []string{v.Addr})
		}
	}

	// 开始读ready
	go rc.serveChannels()
}

// 监听ready通道，集群模式下会开始监听
func (rc *Node) serveChannels() {
	sn, err := rc.raftStorage.Snapshot()
	if err != nil {
		log.Panic(err)
	}
	// 获取当前的日志最新的状态信息，index和term等
	rc.confState = sn.Metadata.ConfState
	rc.snapshotIndex = sn.Metadata.Index
	rc.appliedIndex = sn.Metadata.Index

	defer rc.wal.Close()
	// 心跳定时器
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rc.node.Tick()
		case rd := <-rc.node.Ready():
			if rd.SoftState != nil {
				rc.lead = rd.SoftState.Lead
			}
			//islead = rd.RaftState == raft.StateLeader
			rc.saveToStorage(rd)

			// 将信息发给其他节点
			rc.transport.Send(rd.Messages)

			// 处理需要commit的日志
			ok := rc.publishEntries(rc.entriesToApply(rd.CommittedEntries))
			if !ok {
				//// 此时节点停止
				rc.stop()
				rc.clearConfig()
				return
			}
			// 必要时生成快照
			rc.maybeTriggerSnapshot()
			// 告知底层raft该Ready已处理完
			rc.node.Advance()
		// transport出错
		case err = <-rc.transport.ErrorC:
			log.Info(err.Error())
			rc.stop()
			return
		case <-rc.stopc:
			return
		}
	}
}

// 停止服务相关(暂时不直接关闭程序)
// leave closes http and stops raft.
func (rc *Node) stop() {
	close(rc.stopc)
	rc.transport.Stop()
	rc.node.Stop()
}

// Send 客户端发送propose请求的处理
// 由客户端API调用或Leader收到转发后调用
// 如果是非集群模式(isCluster为false)，直接处理(即service.ProcessHandler后直接service.CommitHandler)
// 如果是集群模式，分两种情况
// 1、当前节点是leader，经service.ProcessHandler后由node.Propose处理后返回，
// 后续会由各个节点的node.Ready读取后进行Commit时由service.CommitHandler处理
// 2、当前节点不是leader，获取当前leader节点地址，转发至leader进行处理(rc.proposeHandler)，
// leader收到请求后经service.ProcessHandler后由node.Propose处理后返回，
// 后续会由各个节点的node.Ready读取后进行Commit时由service.CommitHandler处理
func (rc *Node) Send(namespace, command string, msg interface{}) error {
	// 移除节点后，因为有外部api，故不会停止程序，以此做隔离
	if !rc.active {
		return fmt.Errorf("current node is leave")
	}
	// 非集群模式下直接处理
	if !rc.isCluster {
		data, err := rc.service.ProcessHandler(namespace, command, msg)
		if err != nil {
			return err
		}
		return rc.service.CommitHandler(namespace, data)
	}
	// 集群模式下的处理
	node, isLeader, err := rc.getLeader()
	if err != nil {
		return err
	}
	log.Infof("msg:leader is node(%d)", rc.lead)
	// 如果自己本身就是leader，直接处理，否则转发由leader处理
	if isLeader {
		// service.ProcessHandler要么leader执行，要么非集群模式下自己执行
		data, err := rc.service.ProcessHandler(namespace, command, msg)
		if err != nil {
			return err
		}
		m := &Message{
			From: rc.nodeID,
			Type: PROPOSE,
			Cmd:  namespace,
			Data: data,
		}
		b, err := m.Encode()
		if err != nil {
			return err
		}
		return rc.node.Propose(context.TODO(), b)
	} else {
		cmd, data, err := rc.service.PreProcessHandler(namespace, command, msg)
		if err != nil {
			return err
		}
		return rc.postMessage(node.Addr, cmd, data)
	}
}

// GetPeers 获取集群的peer列表，供API调用
func (rc *Node) GetPeers() (map[uint64]*NodeInfo, uint64, error) {
	if !rc.active {
		return nil, 0, fmt.Errorf("current node is leave")
	}
	return rc.peers.GetAllPeers(), rc.peers.Index(), nil
}

// AddNode 客户端发送增加节点的发送处理
func (rc *Node) AddNode(nodeID uint64, data []byte) error {
	if !rc.active {
		return fmt.Errorf("current node is leave")
	}
	if !rc.isCluster {
		return fmt.Errorf("current node is not cluster mode")
	}
	p := rc.transport.Get(types.ID(nodeID))
	if p != nil {
		return nil
	}
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeID,
		Context: data,
	}
	return rc.node.ProposeConfChange(context.TODO(), cc)
}

// DeleteConfigChange 客户端发送删除节点的发送处理
func (rc *Node) DeleteConfigChange() error {
	if !rc.active {
		return fmt.Errorf("current node is leave")
	}
	if !rc.isCluster {
		return fmt.Errorf("current node is not cluster mode")
	}

	cc := raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: rc.nodeID,
		ID:     rc.nodeID,
	}
	err := rc.node.ProposeConfChange(context.TODO(), cc)
	if err != nil {
		return err
	}
	//rc.leave()
	return nil
}

func (rc *Node) Stop() {
	rc.stop()
}

// InitSend 切换集群时调用，一般一个集群仅调用一次
// 将service现有的缓存信息(基于service.GetInit获取)加载到日志中，便于其他节点同步
// 此时节点刚切换到集群状态，一般会是日志中的第一条信息
// 并通过rc.waiter等待service.ProcessInit处理完后进行后续操作(同步等待)
func (rc *Node) InitSend() error {
	// 集群模式初始化的时候才会调
	if !rc.isCluster {
		return fmt.Errorf("need to change cluster mode")
	}
	cmd, data, err := rc.service.GetInit()
	log.Info("nodeID is ", rc.nodeID)
	if err != nil {
		return err
	}
	m := &Message{
		From: rc.nodeID,
		Type: INIT,
		Cmd:  cmd,
		Data: data,
	}
	b, err := m.Encode()
	if err != nil {
		return err
	}
	err = rc.node.Propose(context.TODO(), b)
	if err != nil {
		log.Error(err)
		return err
	}
	// 等待处理完
	c := rc.waiter.Register(rc.nodeID)
	res := <-c
	str, ok := res.(string)
	if !ok {
		return fmt.Errorf("init send wait channel interface assert error")
	}
	if len(str) > 0 {
		return fmt.Errorf(str)
	}
	return nil
}

// changeCluster 切换集群时调用，一般是非集群节点收到其他节点的join请求时触发(rc.joinHandler)
// 开始运行集群节点,新建日志文件，启动transport和node，
// 并开始监听node.ready,将现有缓存加入日志中rc.InitSend
func (rc *Node) changeCluster(addr string) error {
	log.Info("change cluster mode")
	rc.nodeID = 1
	rc.nodeKey = uuid.New()
	rc.join = true
	rc.isCluster = true
	rc.active = true
	rc.waldir = fmt.Sprintf("eosc-%d", rc.nodeID)
	rc.snapdir = fmt.Sprintf("eosc-%d-snap", rc.nodeID)

	rc.transport.ID = types.ID(rc.nodeID)
	rc.transport.Raft = rc
	rc.transport.LeaderStats = stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(rc.nodeID)))
	rc.transportHandler = rc.genHandler()
	rc.stopc = make(chan struct{})
	// 判断快照文件夹是否存在，不存在则创建
	if !fileutil.Exist(rc.snapdir) {
		if err := os.Mkdir(rc.snapdir, 0750); err != nil {
			return fmt.Errorf("eosc: node(%d) cannot create dir for snapshot (%v)", rc.nodeID, err)
		}
	}
	u, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("eosc: fail to parse address,%w", err)
	}
	rc.broadcastIP = u.Host
	index := strings.Index(u.Host, ":")
	if index > 0 {
		rc.broadcastIP = u.Host[:index]
		rc.broadcastPort, _ = strconv.Atoi(u.Host[index+1:])
	}

	rc.peers.SetPeer(rc.nodeID, &NodeInfo{
		NodeSecret: &NodeSecret{
			ID:  rc.nodeID,
			Key: rc.nodeKey,
		},
		Protocol:      u.Scheme,
		BroadcastIP:   rc.broadcastIP,
		BroadcastPort: rc.broadcastPort,
		Addr:          addr,
	})

	// 新建快照管理
	rc.snapshotter = snap.New(zap.NewExample(), rc.snapdir)

	// 将日志wal重写入raftNode实例中，读取快照和日志，并初始化raftStorage,此处主要是新建日志文件
	rc.wal = rc.replayWAL()
	// 节点配置
	c := &raft.Config{
		ID:                        rc.nodeID,
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   rc.raftStorage,
		MaxSizePerMsg:             1024 * 1024,
		MaxInflightMsgs:           256,
		MaxUncommittedEntriesSize: 1 << 30,
	}
	// 启动节点时添加peers，一般情况下此时只有自己
	peers := make([]raft.Peer, 0, rc.peers.GetPeerNum())
	peerList := rc.peers.GetAllPeers()
	for id, p := range peerList {
		peers = append(peers, raft.Peer{ID: id, Context: p.Marshal()})
	}
	// 启动node节点,新开一个集群
	rc.node = raft.StartNode(c, peers)
	// 通信实例开始运行
	err = rc.transport.Start()
	if err != nil {
		return err
	}
	// 与集群中的其他节点建立通信
	for k, v := range peerList {
		// transport加入peer列表，节点本身不添加
		if k != rc.nodeID {
			rc.transport.AddPeer(types.ID(k), []string{v.Addr})
		}
	}
	// 读ready
	go rc.serveChannels()
	log.Info("change cluster mode successfully")
	go func() {
		// 开始打包处理初始化信息
		err = rc.InitSend()
		if err != nil {
			log.Error(err)
			return
		}
	}()
	rc.writeConfig()
	return nil
}

// 工具方法
// postMessage 转发消息，基于json
func (rc *Node) postMessage(addr string, command string, data []byte) error {
	// 转给leader
	m := &ProposeMsg{
		From: rc.nodeID,
		To:   int(rc.lead),
		Cmd:  command,
		Data: data,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/raft/node/propose", addr), "application/json;charset=utf-8", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res := &Response{}
	err = json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	if res.Code == "000000" {
		return nil
	} else {
		msg := res.Msg
		return fmt.Errorf(msg)
	}
}

// getLeader 获取leader地址以及判断当前节点是不是leader
func (rc *Node) getLeader() (*NodeInfo, bool, error) {
	if rc.lead == raft.None {
		if rc.node.Status().Lead == raft.None {
			return nil, false, fmt.Errorf("current node(%d) has no leader", rc.lead)
		} else {
			rc.lead = rc.node.Status().Lead
		}
	}
	flag := false
	if rc.nodeID == rc.lead {
		flag = true
	}
	v, ok := rc.peers.GetPeerByID(rc.lead)
	if !ok {
		return nil, flag, fmt.Errorf("current node has no leader(%d) host", rc.lead)
	}
	return v, flag, nil
}
