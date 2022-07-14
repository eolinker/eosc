package etcd

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
