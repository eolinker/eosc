package etcd

import (
	"context"
	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/version"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var _ Etcd = (*_Server)(nil)

type _Server struct {
	config                  Config
	ctx                     context.Context
	cancel                  context.CancelFunc
	mu                      sync.RWMutex
	server                  *etcdserver.EtcdServer
	raftHandler             atomic.Pointer[http.Handler]
	leaseHandler            atomic.Pointer[http.Handler]
	downgradeEnabledHandler atomic.Pointer[http.Handler]
	hashKVHandler           atomic.Pointer[http.Handler]
	client                  *clientv3.Client
	requestTimeout          time.Duration
	name                    string
	leaderChangeHandler     []ILeaderStateHandler
	clientCh                []chan *clientv3.Client
	clusterData             *Clusters
}

func (s *_Server) Status() ClusterInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.server == nil {
		return ClusterInfo{}
	}
	member := s.server.Cluster().Member(s.server.ID())
	if member == nil {
		return ClusterInfo{}
	}
	members := s.server.Cluster().Members()
	leaderId := s.server.Leader()

	nodes := s.clusterData.parse(leaderId, members...)

	return ClusterInfo{
		Cluster: s.clusterData.Cluster(),
		Nodes:   nodes,
	}
}

func (s *_Server) Nodes() []*Node {

	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.server != nil {

		members := s.server.Cluster().Members()
		leaderId := s.server.Leader()

		return s.clusterData.parse(leaderId, members...)
	}
	return []*Node{}
}

func (s *_Server) Version() Versions {

	s.mu.RLock()
	defer s.mu.RUnlock()
	strv := "not_decided"
	if s.server != nil {
		v := s.server.Cluster().Version()
		if v != nil {
			strv = v.String()
		}
	}
	return Versions{
		Server:  version.Version,
		Cluster: strv,
	}
}

func NewServer(ctx context.Context, mux *http.ServeMux, config Config) (*_Server, error) {
	serverCtc, cancel := context.WithCancel(ctx)
	s := &_Server{
		config:         config,
		ctx:            serverCtc,
		cancel:         cancel,
		requestTimeout: 10 * time.Second,
		clientCh:       make([]chan *clientv3.Client, 0, 10),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addHandler(mux)

	err := s.initEtcdServer()
	if err != nil {
		return nil, err
	}

	return s, nil
}
func (s *_Server) Info() Info {

	if s.server == nil {
		return nil
	}
	return s.server.Cluster().Member(s.server.ID())
}
func (s *_Server) IsLeader() (bool, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isLeader()
}
func (s *_Server) isLeader() (bool, []string) {
	server := s.server
	lead := server.Leader()
	if lead == server.ID() {
		return true, nil
	}

	return false, server.Cluster().Member(lead).PeerURLs
}

func (s *_Server) Put(key string, value []byte) error {

	ctx, _ := s.requestContext()
	_, err := s.client.Put(ctx, key, string(value))

	return err

}

func (s *_Server) Delete(key string) error {

	ctx, _ := s.requestContext()
	_, err := s.client.Delete(ctx, key)

	return err
}

func (s *_Server) Watch(prefix string, handler ServiceHandler) {
	clientCh := make(chan *clientv3.Client, 1)
	s.mu.Lock()
	s.clientCh = append(s.clientCh, clientCh)
	clientCh <- s.client
	s.mu.Unlock()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		once := sync.Once{}

		var watch clientv3.WatchChan = nil
		for {
			select {
			case client, ok := <-clientCh:
				{
					if !ok {
						continue
					}
					ctx, _ := s.requestContext()
					response, err := client.Get(ctx, prefix, clientv3.WithPrefix())
					if err != nil {
						log.Warn("watch ", prefix, " error:", err)
						return
					}
					watch = client.Watch(s.ctx, prefix, clientv3.WithPrefix())
					init := make([]*KValue, 0, response.Count)
					for _, kv := range response.Kvs {
						init = append(init, &KValue{
							Key:   kv.Key,
							Value: kv.Value,
						})
					}
					handler.Reset(init)
					once.Do(func() {
						wg.Done()
					})
				}

			case <-s.ctx.Done():
				return
			case v, ok := <-watch:
				if !ok {
					watch = nil
					continue
				}
				for _, e := range v.Events {
					switch e.Type {
					case mvccpb.DELETE:
						handler.Delete(string(e.Kv.Key))
					case mvccpb.PUT:
						handler.Put(string(e.Kv.Key), e.Kv.Value)

					}
				}
			}
		}
	}()
	wg.Wait()
}

func (s *_Server) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(s.ctx, s.requestTimeout)
}
