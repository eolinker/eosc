package etcdRaft

import (
	"go.etcd.io/etcd/server/v3/embed"
)

var (
	defaultName               = "apintoNode"
	defaultAutoCompactionMode = embed.CompactorModeRevision
	defaultClientUrls         = []string{"http://127.0.0.1:9400"}
	defaultPeerUrls           = []string{"http://127.0.0.1:9400"}
	defaultClusterMembers     = map[string][]string{defaultName: defaultPeerUrls}
)

const (
	// addr
	defaultDir        = "./data"
	quotaBackendBytes = 8 * 1024 * 1024 * 1024
	maxTxnOps         = 10240
	maxRequestBytes   = 10 * 1024 * 1024 // 10MB
	snapshotCount     = 5000

	// cluster key
	dataPrefixKey      = "/apinto/data"
	defaultClusterName = "APINTO_CLUSTER"
)
