package etcd

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/eolinker/eosc/log"
	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/pkg/v3/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

var (
	_clusterId = []byte("~/cluster")
	_nodePre   = []byte("~/nodes/")
)

type NodeGatewayConfig struct {
	Urls []string `json:"urls"`
}
type Clusters struct {
	locker  sync.RWMutex
	data    map[string]*NodeGatewayConfig
	cluster string

	client *clientv3.Client
}

func (cs *Clusters) Cluster() string {
	cs.locker.RLock()
	defer cs.locker.RUnlock()
	return cs.cluster
}

func (cs *Clusters) SetCluster(cluster string) {
	cs.locker.Lock()
	defer cs.locker.Unlock()
	cs.cluster = cluster

	cs.client.Put(cs.client.Ctx(), string(_clusterId), cluster)
}

type EventType = mvccpb.Event_EventType

func NewClusters(ctx context.Context, client *clientv3.Client) *Clusters {
	c := &Clusters{
		locker:  sync.RWMutex{},
		cluster: "",
		data:    map[string]*NodeGatewayConfig{},
		client:  client,
	}

	response, err := client.Get(ctx, "~/", clientv3.WithPrefix())
	if err != nil {
		log.Warn("get init cluster:", err)
		return c
	}

	watch := client.Watch(ctx, "~/", clientv3.WithPrefix())

	for _, kv := range response.Kvs {
		c.doEvent(mvccpb.PUT, kv.Key, kv.Value)
	}
	if c.cluster == "" {
		c.cluster = uuid.NewString()
		client.Put(ctx, string(_clusterId), c.cluster)
	}
	go func() {
		for watcher := range watch {
			c.locker.Lock()
			for _, event := range watcher.Events {
				c.doEvent(event.Type, event.Kv.Key, event.Kv.Value)
			}
			c.locker.Unlock()
		}
	}()
	return c
}
func (cs *Clusters) doEvent(t EventType, key, v []byte) {
	if bytes.Equal(key, _clusterId) {
		cs.cluster = string(v)
		return
	}

	nodeId := string(bytes.TrimPrefix(key, _nodePre))

	switch t {
	case mvccpb.PUT:
		config := new(NodeGatewayConfig)
		json.Unmarshal(v, config)
		cs.data[nodeId] = config
	case mvccpb.DELETE:
		delete(cs.data, nodeId)
	}
}

func (cs *Clusters) parse(leader types.ID, members ...Info) []*Node {
	nodes := make([]*Node, 0, len(members))
	cs.locker.RLock()
	defer cs.locker.RUnlock()
	for _, m := range members {
		n := &Node{
			Id:       m.ID.String(),
			Name:     m.Name,
			Peer:     m.PeerURLs,
			Admin:    m.ClientURLs,
			IsLeader: leader == m.ID,
		}
		if g, has := cs.data[n.Id]; has {
			n.Server = g.Urls
		}
		nodes = append(nodes, n)
	}
	return nodes
}
