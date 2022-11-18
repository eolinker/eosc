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
	data    map[string]*NodeGatewayConfig
	cluster string
	mu      sync.RWMutex
	client  *clientv3.Client
}

type EventType = mvccpb.Event_EventType

func NewClusters(ctx context.Context, client *clientv3.Client, s *_Server) *Clusters {
	c := &Clusters{
		cluster: "",
		data:    map[string]*NodeGatewayConfig{},
	}

	response, err := client.Get(ctx, "~/", clientv3.WithPrefix())
	if err != nil {
		log.Warn("get init cluster:", err)
		return c
	}

	watch := client.Watch(ctx, "~/", clientv3.WithPrefix())

	for _, kv := range response.Kvs {
		if bytes.Equal(kv.Key, _clusterId) {
			c.cluster = string(kv.Value)
			continue
		}
		nodeId := string(bytes.TrimPrefix(kv.Key, _nodePre))
		config := new(NodeGatewayConfig)
		json.Unmarshal(kv.Value, config)
		c.data[nodeId] = config
	}
	if c.cluster == "" {
		c.cluster = uuid.NewString()
		client.Put(ctx, string(_clusterId), c.cluster)
	}
	go func() {
		for watcher := range watch {
			c.mu.Lock()
			for _, event := range watcher.Events {
				c.nodeEventDoer(event.Type, event.Kv.Key, event.Kv.Value)
			}

			memberInitUrls := make(map[string][]string)
			for _, m := range s.server.Cluster().Members() {
				id := m.Name
				if _, has := c.data[id]; has {
					memberInitUrls[id] = m.PeerURLs
				}
			}
			clusterString := initialClusterString(memberInitUrls)
			s.resetCluster(clusterString)
			c.mu.Unlock()
		}
	}()
	return c
}
func (cs *Clusters) nodeEventDoer(t EventType, key, v []byte) {
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
