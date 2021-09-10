package raft

import (
	"net/http"
	"testing"

	"github.com/eolinker/eosc/log"

	raft_service "github.com/eolinker/eosc/raft/raft-service"
	store2 "github.com/eolinker/eosc/store"
)

func TestRaftNode1(t *testing.T) {
	store, _ := store2.NewStore()
	node := NewNode(raft_service.NewService(store))
	sm := http.NewServeMux()

	sm.Handle("/raft/", node)
	sm.Handle("/raft", node)

	log.Fatal(http.ListenAndServe(":9999", sm))
}

func TestRaftNode2(t *testing.T) {
	store, _ := store2.NewStore()
	service := raft_service.NewService(store)
	node := NewNode(service)
	err := JoinCluster(node, "127.0.0.1", 9998, "http://127.0.0.1:9999", "http://127.0.0.1:9999/raft/node/join", "http", service, 0)
	if err != nil {
		log.Error(err)
		return
	}
	sm := http.NewServeMux()
	sm.Handle("/raft/", node)
	sm.Handle("/raft", node)
	log.Fatal(http.ListenAndServe(":9998", sm))
}

func TestRaftNode3(t *testing.T) {
	store, _ := store2.NewStore()
	service := raft_service.NewService(store)
	node := NewNode(service)
	err := JoinCluster(node, "127.0.0.1", 9997, "http://127.0.0.1:9999", "http://127.0.0.1:9999/raft/node/join", "http", service, 0)
	if err != nil {
		log.Error(err)
		return
	}
	sm := http.NewServeMux()
	sm.Handle("/raft/", node)
	sm.Handle("/raft", node)
	log.Fatal(http.ListenAndServe(":9997", sm))
}
