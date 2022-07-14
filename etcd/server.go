package etcd

import (
	"context"
	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"net/http"
	"sync"
	"time"
)

var _ Etcd = (*_Server)(nil)

type _Server struct {
	ctx                     context.Context
	cancel                  context.CancelFunc
	mu                      sync.RWMutex
	server                  *etcdserver.EtcdServer
	raftHandler             http.Handler
	leaseHandler            http.Handler
	downgradeEnabledHandler http.Handler
	hashKVHandler           http.Handler
	client                  *clientv3.Client
	requestTimeout          time.Duration
	name                    string
	leaderChangeHandler     []ILeaderStateHandler
}

func NewServer(ctx context.Context, mux *http.ServeMux) (*_Server, error) {
	serverCtc, cancel := context.WithCancel(ctx)
	s := &_Server{
		ctx:            serverCtc,
		cancel:         cancel,
		requestTimeout: 10 * time.Second,
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
	s.mu.RLock()
	defer s.mu.RUnlock()
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

	s.mu.RLock()
	defer s.mu.RUnlock()
	ctx, _ := s.requestContext()
	_, err := s.client.Put(ctx, key, string(value))

	return err

}

func (s *_Server) Delete(key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ctx, _ := s.requestContext()
	_, err := s.client.Delete(ctx, key)

	return err
}

func (s *_Server) Watch(prefix string, handler ServiceHandler) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ctx, _ := s.requestContext()
	response, err := s.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Warn("watch ", prefix, " error:", err)
		return
	}

	watch := s.client.Watch(s.ctx, prefix, clientv3.WithPrefix())

	go func() {

		init := make([]*KValue, 0, response.Count)
		for _, kv := range response.Kvs {
			init = append(init, &KValue{
				Key:   kv.Key,
				Value: kv.Value,
			})
		}
		handler.Reset(init)
		for {
			select {
			case <-s.ctx.Done():
				return
			case v, ok := <-watch:
				if !ok {
					return
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
}

func (s *_Server) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(s.ctx, s.requestTimeout)
}
