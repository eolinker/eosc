package etcdRaft

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
	"log"
	"net/http"
	"sync"
	"time"
)

var _ Etcd = (*EtcdServer)(nil)

type EtcdServer struct {
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	server         *etcdserver.EtcdServer
	client         *clientv3.Client
	requestTimeout time.Duration
	raftHandler    http.Handler
}

func NewEtcdNode(name string, clients []string, peers []string, clusters map[string][]string) (*EtcdServer, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &EtcdServer{
		ctx:            ctx,
		cancel:         cancelFunc,
		requestTimeout: 10 * time.Second,
	}
	err := s.start(name, clients, peers, clusters)
	return s, err
}

func (e *EtcdServer) start(name string, clients []string, peers []string, clusters map[string][]string) error {
	server, err := NewEtcdServer(name, clients, peers, clusters)
	if err != nil {
		return err
	}
	// block: wait for server ready
	select {
	case <-server.ReadyNotify():
		log.Print("Server is ready!")
		e.server = server
	case <-time.After(60 * time.Second):
		closeEtcdServer(server) // trigger a shutdown
		return errors.New("server took too long to start!" +
			"start others if they're not online, " +
			"otherwise purge this member, clean data directory " +
			"and rejoin it back")
	}
	e.raftHandler = e.genHandler()
	e.client = v3client.New(server)
	return nil
}

// Restart 加入新集群和离开集群的时候需要restart server， join操作需要清空缓存，leave操作需要reset缓存
func (e *EtcdServer) restart(name string, peers []string, clients []string, clusters map[string][]string) error {
	err := e.close()
	if err != nil {
		return err
	}
	// 清楚旧的日志文件
	err = e.cleanWalFile()
	if err != nil {
		return err
	}
	return e.start(name, peers, clients, clusters)

}

func (e *EtcdServer) Join() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 	发出join请求，获取相关信息

	//e.restart()

	return nil
}

func (e *EtcdServer) Leave() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	// 获取全部数据
	allData, err := e.getAllData()
	if err != nil {
		return err
	}
	// 当前节点
	current := e.server.Cluster().Member(e.server.ID())
	cluster := map[string][]string{
		current.Name: current.PeerURLs,
	}
	// todo 写本地缓存操作 writeConfig

	// leave相关操作
	// 集群中删除自己
	err = e.removeMember(uint64(current.ID))
	if err != nil {
		return err
	}
	// 重启etcd服务
	err = e.restart(current.Name, current.PeerURLs, current.ClientURLs, cluster)
	if err != nil {
		return err
	}
	e.resetAllData(allData)
	return nil
}

func (e *EtcdServer) IsLeader() (bool, []string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	server := e.server
	lead := server.Leader()
	if lead == server.ID() {
		return true, nil, nil
	}

	return false, server.Cluster().Member(lead).PeerURLs, nil
}

// Close 关闭集群
func (e *EtcdServer) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cancel()
	return e.close()
}
func (e *EtcdServer) close() error {

	// 关闭客户端
	e.raftHandler = nil
	e.closeClient()
	// 关闭etcd server
	e.closeServer()
	return nil
}

// 关闭客户端
func (e *EtcdServer) closeClient() {

	if e.client == nil {
		return
	}
	err := e.client.Close()
	if err != nil {
		log.Printf("close client failed: %v", err)
	}
	e.client = nil
}

// 关闭etcd服务
func closeEtcdServer(s *etcdserver.EtcdServer) {
	select {
	case <-s.ReadyNotify():
		s.Stop()
		<-s.StopNotify()
	default:
		s.HardStop()
		log.Printf("hard stop server")
	}
}

// 停止服务
func (e *EtcdServer) closeServer() {

	if e.server == nil {
		return
	}
	closeEtcdServer(e.server)
	e.server = nil
}
