package etcd

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/pkg/v3/types"
	clientv3 "go.etcd.io/etcd/client/v3"
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
}

type EventType = mvccpb.Event_EventType

func getClusterId() string {
	clusterId, has := env.GetEnv("CLUSTER_ID")
	if !has || clusterId == "" {
		return uuid.NewString()
	}
	return clusterId
}

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
		cfg := new(NodeGatewayConfig)
		_ = json.Unmarshal(kv.Value, cfg)
		c.data[nodeId] = cfg
	}
	if c.cluster == "" {
		c.cluster = getClusterId()
		_, _ = client.Put(ctx, string(_clusterId), c.cluster)
	}
	go func() {
		for watcher := range watch {

			c.mu.Lock()
			for _, event := range watcher.Events {
				//if event.Type == mvccpb.DELETE {
				//	log.DebugF("node event: %s, %s, %s", event.Type, string(event.Kv.Key), s.server.ID().String())
				//	nodeId := string(bytes.TrimPrefix(event.Kv.Key, _nodePre))
				//	if nodeId == s.server.ID().String() {
				//		// while remove self
				//		allData, _ := s.getAllData()
				//		s.clearCluster()
				//		err = s.restart("")
				//		if err != nil {
				//			log.Errorf("restart error: %s", err.Error())
				//			return
				//		}
				//
				//		s.resetAllData(allData)
				//		return
				//	}
				//}
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
	//go func() {
	//	ticket := time.NewTicker(5 * time.Second)
	//	defer ticket.Stop()
	//	var allData map[string][]byte
	//	var err error
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		case <-s.server.StopNotify():
	//			s.clearCluster()
	//
	//			err = s.restart("")
	//			if err != nil {
	//				log.Errorf("restart error: %s", err.Error())
	//				return
	//			}
	//			if allData != nil {
	//				s.resetAllData(allData)
	//			}
	//			return
	//		case <-ticket.C:
	//			allData, err = s.getAllData()
	//			if err != nil {
	//				log.Errorf("get all data error: %s", err.Error())
	//			}
	//		}
	//	}
	//}()
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
		_ = json.Unmarshal(v, config)
		cs.data[nodeId] = config
	case mvccpb.DELETE:
		delete(cs.data, nodeId)
	}
}

func (cs *Clusters) parse(leader types.ID, members ...Info) []*Node {
	nodes := make([]*Node, 0, len(members))

	for _, m := range members {
		n := &Node{
			Id:       m.ID,
			ID:       m.ID.String(),
			Name:     m.Name,
			Peer:     m.PeerURLs,
			Admin:    m.ClientURLs,
			IsLeader: leader == m.ID,
		}
		if g, has := cs.data[n.ID]; has {
			n.Server = g.Urls
		}
		nodes = append(nodes, n)
	}
	return nodes
}
