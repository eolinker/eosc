package main

import (
	"flag"
)

// 节点1 raft.exe --id 1 --keys 1 --cluster http://127.0.0.1:12379 --port 12380
// 节点2 raft.exe --id 2 --keys 1,2 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379 --port 22380 --join true
// 节点3 raft.exe --id 3 --keys 1,2,3 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port 32380 --join true
// 节点4 raft.exe --id 4 --keys 1,2,3,4 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379,http://127.0.0.1:42379 --port 42380 --join true
// 节点5 raft.exe --id 5 --keys 1,2,3,4,5 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379,http://127.0.0.1:42379,http://127.0.0.1:52379 --port 52380 --join true

var (
	httpPort   = 1234
	nodeID = 1
	isCluster = true
	join = false
	host = "http://127.0.0.1:8081"
	keys = "1"
	peers = "http://127.0.0.1:8081"
	target = "http://127.0.0.1:8081"
)

func initFlag() {
	// 本地节点http端口
	flag.IntVar(&httpPort, "http", 8081, "Please provide a valid http port for api")
	// 本地节点ID(集群中唯一)
	flag.IntVar(&nodeID, "id", 1, "node ID")
	// 本地节点地址
	flag.StringVar(&host, "host", "http://127.0.0.1:8081", "node host")
	// 开启集群模式
	flag.BoolVar(&isCluster, "isCluster", true, "cluster mode")
	// peer列表，以,分割，非集群模式下可为空
	flag.StringVar(&peers, "peers", "http://127.0.0.1:8081", "comma separated cluster peers")
	// peer对应的节点id列表，以,分割，非集群模式下可为空
	flag.StringVar(&keys, "keys", "1", "comma separated node ids")
	// 是否加入一个新集群，非集群模式下可为空
	flag.BoolVar(&join, "join", false, "join an existing cluster")
	// 目标集群节点地址，join存在时不能为空
	flag.StringVar(&target, "target", "http://127.0.0.1:8081", "target cluster")
	flag.Parse()
}


