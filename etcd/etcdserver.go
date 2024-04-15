package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/client/pkg/v3/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/config"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
	"go.etcd.io/etcd/server/v3/wal"
)

var (
	ErrorAlreadyInCluster = errors.New("already in cluster")
	ErrorNotInCluster     = errors.New("not in cluster")
	ErrorMemberNotExist   = errors.New("member not exist")
)

func (s *_Server) initEtcdServer() error {

	c := s.etcdServerConfig()
	defer func() {
		if c.NewCluster {
			_, _ = http.Get(fmt.Sprintf("https://statistics.apinto.com/report/deploy/n?id=%s", c.Name))
		}
	}()
	s.name = c.Name
	srv, err := createEtcdServer(c)
	if err != nil {
		return err
	}
	raftHandler := srv.RaftHandler()
	s.raftHandler.Swap(&raftHandler)
	leaseHandler := srv.LeaseHandler()
	s.leaseHandler.Swap(&leaseHandler)
	downgradeEnabledHandler := srv.DowngradeEnabledHandler()
	s.downgradeEnabledHandler.Swap(&downgradeEnabledHandler)
	hashKVHandler := srv.HashKVHandler()
	s.hashKVHandler.Swap(&hashKVHandler)
	s.server = srv
	go s.check(srv)
	<-s.server.ReadyNotify()

	s.client = v3client.New(s.server)
	gatewayConfig := &NodeGatewayConfig{Urls: s.config.GatewayAdvertiseUrls}
	data, _ := json.Marshal(gatewayConfig)
	_, err = s.client.Put(s.ctx, fmt.Sprintf("~/nodes/%s", s.server.ID()), string(data))
	if err != nil {
		return err
	}
	s.clusterData = NewClusters(s.ctx, s.client, s)
	_ = os.Setenv("cluster_id", s.clusterData.cluster)
	for _, ch := range s.clientCh {
		ch <- s.client
	}
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
	if memberInitialized && srvcfg.InitialCorruptCheck {
		if err = server.CorruptionChecker().InitialCheck(); err != nil {
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
func (s *_Server) restart(InitialCluster string) error {
	err := s.close()
	if err != nil {
		return err
	}
	// 清楚旧的日志文件
	err = s.cleanWalFile()
	if err != nil {
		return err
	}
	s.resetCluster(InitialCluster)
	return s.initEtcdServer()
}
func (s *_Server) checkIsJoined() bool {
	etcdInitPath := filepath.Join(s.config.DataDir, "cluster", "etcd.init")

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

	if s.checkIsJoined() {
		return ErrorAlreadyInCluster
	}
	urls, clientUrls, err := s.config.CreatePeerUrl()
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

	return s.restart(InitialCluster)
}

func (s *_Server) Remove(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.server == nil {
		return ErrorNotInCluster
	}
	// 获取全部数据
	currentId := s.server.ID()
	members := s.server.Cluster().Members()
	if len(members) == 1 && members[0].ID == currentId {
		return ErrorNotInCluster
	}
	for _, m := range members {
		if m.Name == name {
			if m.ID == s.server.Leader() {
				return fmt.Errorf("cannot remove leader")
			}
			_, err := s.client.Delete(s.ctx, fmt.Sprintf("~/nodes/%s", m.ID))
			if err != nil {
				return err
			}
			err = s.removeMember(uint64(m.ID))
			if err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("%w name %s", ErrorMemberNotExist, name)
}

func (s *_Server) Leave() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server == nil {
		return ErrorNotInCluster
	}

	// 获取全部数据
	currentId := s.server.ID()
	members := s.server.Cluster().Members()
	if len(members) == 1 && members[0].ID == currentId {
		return ErrorNotInCluster
	}
	allData, err := s.getAllData()
	if err != nil {
		return err
	}
	if err = s.server.TransferLeadership(); err != nil {
		log.Warn("leadership transfer failed ", currentId, " error:", err)
		return err
	}
	current := s.server.Cluster().Member(currentId)

	// leave相关操作
	// 集群中删除自己
	_, err = s.client.Delete(s.ctx, fmt.Sprintf("~/nodes/%s", s.server.ID()))
	if err != nil {
		return err
	}

	err = s.removeMember(uint64(current.ID))
	if err != nil {
		return err
	}
	// 清楚集群配置
	s.clearCluster()
	// 重启etcd服务

	err = s.restart("")
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
	emptyHandler := http.NotFoundHandler()
	s.raftHandler.Swap(&emptyHandler)
	s.hashKVHandler.Swap(&emptyHandler)
	s.downgradeEnabledHandler.Swap(&emptyHandler)
	s.leaseHandler.Swap(&emptyHandler)
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

}
