package etcdRaft

import (
	"context"
	"errors"

	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
	"log"
	"net"
	"net/http"
	"net/url"
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
	peers []net.Listener
}

func NewEtcdNode(name string, clients []string, peers []string, clusters map[string][]string) (*EtcdServer, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &EtcdServer{
		ctx:            ctx,
		cancel:         cancelFunc,
		requestTimeout: 10 * time.Second,
		peers: make([]net.Listener, 0),
	}
	err := s.start(name, clients, peers, clusters, false)
	return s, err
}

func (e *EtcdServer) start(name string, clients []string, peers []string, clusters map[string][]string, isJoin bool) (err error) {
	e.server, err = NewEtcdServer(name, clients, peers, clusters, isJoin)
	if err != nil {
		return err
	}
	e.raftHandler = e.genHandler()

	// Todo 必须在server.Start()后监听和handler,这样节点间才能正常通信，etcd服务才算启动成功（所以需要在server.ReadyNotify()前调用监听），看怎么跟外部业务整合（路由树那里）
	// Todo 因为可能跟外部的监听器整合，此处暂时存储亿方便listener的close逻辑，需注意
	for _, peer := range peers {
		u, err := url.Parse(peer)
		listen, err := transport.NewListenerWithOpts(u.Host, u.Scheme,
			transport.WithTimeout(rafthttp.ConnReadTimeout, rafthttp.ConnWriteTimeout),
		)
		if err != nil {
			log.Print(err)
			continue
		}
		srv := &http.Server{
			Handler:     e.raftHandler,
			ReadTimeout: 5 * time.Minute,
		}
		go srv.Serve(listen)
		e.peers = append(e.peers, listen)
	}

	// block: wait for server ready
	select {
	// Todo 这里是等到监听成功才赋值，目的是确保server启动成功，看有没有必要保留（考虑到http handler需要等到路由树注册的时候才监听，可以考虑这里的处理直接不要）
	case <-e.server.ReadyNotify():
		log.Print("Server is ready!")
	case <-time.After(60 * time.Second):
		closeEtcdServer(e.server) // trigger a shutdown
		e.server = nil
		return errors.New("server took too long to start!" +
			"start others if they're not online, " +
			"otherwise purge this member, clean data directory " +
			"and rejoin it back")
	}

	e.client = v3client.New(e.server)
	return nil
}

// Restart 加入新集群和离开集群的时候需要restart server， join操作需要清空缓存，leave操作需要reset缓存
func (e *EtcdServer) restart(name string, peers []string, clients []string, clusters map[string][]string, isJoin bool) error {
	err := e.close()
	if err != nil {
		return err
	}
	// 清楚旧的日志文件
	err = e.cleanWalFile()
	if err != nil {
		return err
	}
	return e.start(name, peers, clients, clusters, isJoin)

}

func (e *EtcdServer) Join(target string, addr []string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 	发出join请求，获取相关信息
	name, clusters, err := e.sendJoinRequest(target, addr)
	if err != nil {
		return err
	}


	return e.restart(name, clusters[name], clusters[name], clusters, true)
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
	e.server.Cluster().ID()
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
	err = e.restart(current.Name, current.PeerURLs, current.ClientURLs, cluster, false)
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
	// todo 临时在这里close
	for _, peer := range e.peers {
		_ = peer.Close()
	}
}
