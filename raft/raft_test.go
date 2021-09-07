package raft

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/eolinker/eosc/log"
)

func TestRaft(t *testing.T) {
	initFlag()

	// 初始化服务
	var s = Create()
	var raft = &Node{}
	var err error
	if !join {
		// 新建raft节点,以集群模式启动或非集群单点模式
		raft, err = CreateRaftNode(nodeID, host, s, peers, keys, join, isCluster)
	} else {
		// 新建raft节点,加入一个集群
		raft, err = JoinCluster(host, target, s)
	}
	if err != nil {
		log.Fatal(err)
	}
	client := &Client{
		raft: raft,
	}

	//httpServer := http.NewServeMux()
	//httpServer.Handle("/raft/api/", client.Handler())
	log.Info(fmt.Sprintf("Listen http port %d successfully", httpPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), client.Handler()))
}
