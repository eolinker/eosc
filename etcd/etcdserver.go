package etcd

import (
	"errors"
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/client/pkg/v3/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/config"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
	"go.etcd.io/etcd/server/v3/wal"
	"path/filepath"
)

var (
	ErrorAlreadyInCluster = errors.New("already in cluster")
)

func (s *_Server) initEtcdServer() error {

	c := etcdServerConfig()
	s.name = c.Name
	srv, err := createEtcdServer(c)
	if err != nil {
		return err
	}

	s.raftHandler = srv.RaftHandler()
	s.leaseHandler = srv.LeaseHandler()
	s.downgradeEnabledHandler = srv.DowngradeEnabledHandler()
	s.hashKVHandler = srv.HashKVHandler()
	s.server = srv
	go s.check(srv)
	<-s.server.ReadyNotify()
	s.client = v3client.New(s.server)

	return nil

}
func (s *_Server) check(srv *etcdserver.EtcdServer) {
	log.Debug("start check LeaderChanged")
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-srv.LeaderChangedNotify():
			{
				log.Debug("Leader changed")
				isLeader, _ := s.IsLeader()
				hs := s.getLeaderChangeHandlers()
				for _, h := range hs {
					h.LeaderChange(isLeader)
				}
			}

		}
	}
}
func (s *_Server) getLeaderChangeHandlers() []ILeaderStateHandler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hs := s.leaderChangeHandler
	return hs
}
func (s *_Server) HandlerLeader(hs ...ILeaderStateHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	isLeader, _ := s.isLeader()
	for _, h := range hs {
		h.LeaderChange(isLeader)
	}
	s.leaderChangeHandler = append(s.leaderChangeHandler, hs...)

}
func createEtcdServer(srvcfg config.ServerConfig) (*etcdserver.EtcdServer, error) {
	memberInitialized := wal.Exist(filepath.Join(srvcfg.DataDir, "member", "wal"))
	server, err := etcdserver.NewServer(srvcfg)
	if err != nil {
		return nil, err
	}
	if memberInitialized {
		if err = server.CheckInitialHashKV(); err != nil {
			log.Warn("checkInitialHashKV failed", err)
			server.Cleanup()
			server = nil
			return nil, err
		}
	}
	server.Start()
	return server, nil
}

// Restart 加入新集群和离开集群的时候需要restart server， join操作需要清空缓存，leave操作需要reset缓存
func (s *_Server) restart() error {
	err := s.close()
	if err != nil {
		return err
	}
	// 清楚旧的日志文件
	err = s.cleanWalFile()
	if err != nil {
		return err
	}
	return s.initEtcdServer()
}
func checkIsJoined() bool {
	etcdConfig := env.NewConfig(etcdInitPath)
	etcdConfig.ReadFile(etcdInitPath)
	InitialCluster, has := etcdConfig.Get("cluster")
	if !has {
		return false
	}

	urlsMap, err := types.NewURLsMap(InitialCluster)
	if err != nil {
		return false
	}
	if len(urlsMap) > 1 {
		return true
	}
	return false
}
func (s *_Server) Join(target string) error {

	if checkIsJoined() {
		return ErrorAlreadyInCluster
	}
	urls, clientUrls, err := CreatePeerUrl()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// 	发出join请求，获取相关信息
	clusters, err := s.sendJoinRequest(target, urls.StringSlice(), clientUrls.StringSlice())
	if err != nil {
		return err
	}
	InitialCluster := initialClusterString(clusters)

	resetCluster(InitialCluster)

	return s.restart()
}

func (s *_Server) Leave() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 获取全部数据
	allData, err := s.getAllData()
	if err != nil {
		return err
	}
	// 当前节点
	s.server.Cluster().ID()
	current := s.server.Cluster().Member(s.server.ID())

	// leave相关操作
	// 集群中删除自己
	err = s.removeMember(uint64(current.ID))
	if err != nil {
		return err
	}
	// 清楚集群配置
	clearCluster()
	// 重启etcd服务

	err = s.restart()
	if err != nil {
		return err
	}
	s.resetAllData(allData)
	return nil
}

// Close 关闭集群
func (s *_Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancel()
	return s.close()
}
func (s *_Server) close() error {
	s.raftHandler = nil
	s.hashKVHandler = nil
	s.downgradeEnabledHandler = nil
	s.leaseHandler = nil
	s.closeClient()
	// 关闭etcd server
	s.closeServer()
	return nil
}

// 关闭客户端
func (s *_Server) closeClient() {

	if s.client == nil {
		return
	}
	err := s.client.Close()
	if err != nil {
		log.Info("close client failed: %v", err)
	}
	s.client = nil
}

// 关闭etcd服务
func closeEtcdServer(s *etcdserver.EtcdServer) {
	select {
	case <-s.ReadyNotify():
		s.Stop()
		<-s.StopNotify()
	default:
		s.HardStop()
		log.Info("hard stop server")
	}
}

// 停止服务
func (s *_Server) closeServer() {

	if s.server == nil {
		return
	}
	closeEtcdServer(s.server)
	s.server = nil
}

func (s *_Server) getAllData() (map[string][]byte, error) {
	client := s.client
	resp, err := func() (*clientv3.GetResponse, error) {
		ctx, cancel := s.requestContext()
		defer cancel()
		return client.Get(ctx, "/", clientv3.WithPrefix())
	}()
	if err != nil {
		return nil, err
	}
	kvs := make(map[string][]byte)
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = kv.Value
	}
	return kvs, nil
}

func (s *_Server) removeMember(id uint64) error {
	ctx, cancel := s.requestContext()
	defer cancel()
	_, err := s.server.RemoveMember(ctx, id)
	return err
}

func (s *_Server) resetAllData(data map[string][]byte) {
	client := s.client

	for key, bytes := range data {
		_, err := client.Put(s.ctx, key, string(bytes))
		if err != nil {
			log.Warn("reset all data error : %s", err.Error())
		}
	}
	return
}
