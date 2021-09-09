package raft

import (
	"net/http"
	"testing"

	"github.com/eolinker/eosc/log"

	raft_service "github.com/eolinker/eosc/raft/raft-service"
	store2 "github.com/eolinker/eosc/store"
)

//func TestRaft(t *testing.T) {
//	initFlag()
//
//	// 初始化服务
//	var s = Create()
//	var raft = &Node{}
//	var err error
//	if !join {
//		// 新建raft节点,以集群模式启动或非集群单点模式
//		raft, err = CreateRaftNode(nodeID, host, s, peers, keys, join, isCluster)
//	} else {
//		// 新建raft节点,加入一个集群
//		raft, err = JoinCluster(host, target, s)
//	}
//	if err != nil {
//		log.Fatal(err)
//	}
//	client := &Client{
//		raft: raft,
//	}
//
//	//httpServer := http.NewServeMux()
//	//httpServer.Handle("/raft/api/", client.Handler())
//	log.Info(fmt.Sprintf("Listen http port %d successfully", httpPort))
//	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), client.Handler()))
//}

func TestRaftNode1(t *testing.T) {
	store, _ := store2.NewStore()
	node, _ := NewNode(raft_service.NewService(store))
	sm := http.NewServeMux()

	sm.Handle("/raft/node/", node.Handler())
	go func() {
		for {
			select {
			case v, ok := <-node.updateTransport:
				{
					if !ok {
						return
					}
					if v {
						sm.Handle("/", node.transportHandler)
					}
				}
			}
		}
	}()
	log.Fatal(http.ListenAndServe(":9999", sm))
}

func TestRaftNode2(t *testing.T) {
	store, _ := store2.NewStore()
	node, err := JoinCluster("127.0.0.1", 9998, "http://127.0.0.1:9999", raft_service.NewService(store), 0)
	if err != nil {
		log.Error(err)
		return
	}
	sm := http.NewServeMux()
	sm.Handle("/raft/node/", node.Handler())
	sm.Handle("/", node.transport.Handler())
	log.Fatal(http.ListenAndServe(":9998", sm))
}
