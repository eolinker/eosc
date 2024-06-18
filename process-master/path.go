package process_master

import (
	"fmt"
	"github.com/eolinker/eosc/router"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
)

var (
	etcdPaths = []string{
		"/raft/node/join",
		rafthttp.RaftPrefix,
		"/members",
		"/version",
		"/leases",
		"/downgrade/enabled",
	}
	masterApiPaths = []string{
		"/system/",
		fmt.Sprintf("%slog/node/", router.RouterPrefix),
	}
	workerApiPaths = []string{
		router.RouterPrefix,
	}
)
