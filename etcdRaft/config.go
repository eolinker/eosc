package etcdRaft

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/config"
	"go.etcd.io/etcd/server/v3/embed"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/wal"
	"go.uber.org/zap"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func startEtcdServer(cfg *embed.Config) (*etcdserver.EtcdServer, error) {
	var (
		urlsmap types.URLsMap
		token   string
		err     error
	)
	if err = cfg.Validate(); err != nil {
		return nil, err
	}
	memberInitialized := true
	if !isMemberInitialized(cfg) {
		memberInitialized = false
		urlsmap, token, err = cfg.PeerURLsMapAndToken("etcd")
		if err != nil {
			return nil, fmt.Errorf("error setting up initial cluster: %v", err)
		}
	}
	if err != nil {
		return nil, err
	}
	srvcfg := config.ServerConfig{
		Name:                                     cfg.Name,
		ClientURLs:                               cfg.ACUrls,
		PeerURLs:                                 cfg.APUrls,
		DataDir:                                  cfg.Dir,
		DedicatedWALDir:                          cfg.WalDir,
		SnapshotCount:                            snapshotCount,
		SnapshotCatchUpEntries:                   cfg.SnapshotCatchUpEntries,
		MaxSnapFiles:                             cfg.MaxSnapFiles,
		MaxWALFiles:                              cfg.MaxWalFiles,
		InitialPeerURLsMap:                       urlsmap,
		InitialClusterToken:                      token,
		DiscoveryURL:                             cfg.Durl,
		DiscoveryProxy:                           cfg.Dproxy,
		NewCluster:                               cfg.IsNewCluster(),
		PeerTLSInfo:                              cfg.PeerTLSInfo,
		TickMs:                                   cfg.TickMs,
		ElectionTicks:                            cfg.ElectionTicks(),
		InitialElectionTickAdvance:               cfg.InitialElectionTickAdvance,
		AutoCompactionRetention:                  time.Duration(10),
		AutoCompactionMode:                       cfg.AutoCompactionMode,
		QuotaBackendBytes:                        quotaBackendBytes,
		BackendBatchLimit:                        cfg.BackendBatchLimit,
		BackendFreelistType:                      bolt.FreelistMapType,
		BackendBatchInterval:                     cfg.BackendBatchInterval,
		MaxTxnOps:                                maxTxnOps,
		MaxRequestBytes:                          maxRequestBytes,
		SocketOpts:                               transport.SocketOpts{},
		StrictReconfigCheck:                      cfg.StrictReconfigCheck,
		ClientCertAuthEnabled:                    cfg.ClientTLSInfo.ClientCertAuth,
		AuthToken:                                cfg.AuthToken,
		BcryptCost:                               cfg.BcryptCost,
		TokenTTL:                                 cfg.AuthTokenTTL,
		CORS:                                     cfg.CORS,
		HostWhitelist:                            cfg.HostWhitelist,
		InitialCorruptCheck:                      cfg.ExperimentalInitialCorruptCheck,
		CorruptCheckTime:                         cfg.ExperimentalCorruptCheckTime,
		PreVote:                                  cfg.PreVote,
		Logger:                                   cfg.GetLogger(),
		ForceNewCluster:                          cfg.ForceNewCluster,
		EnableGRPCGateway:                        cfg.EnableGRPCGateway,
		ExperimentalEnableDistributedTracing:     cfg.ExperimentalEnableDistributedTracing,
		UnsafeNoFsync:                            cfg.UnsafeNoFsync,
		EnableLeaseCheckpoint:                    cfg.ExperimentalEnableLeaseCheckpoint,
		CompactionBatchLimit:                     cfg.ExperimentalCompactionBatchLimit,
		WatchProgressNotifyInterval:              cfg.ExperimentalWatchProgressNotifyInterval,
		DowngradeCheckTime:                       cfg.ExperimentalDowngradeCheckTime,
		WarningApplyDuration:                     cfg.ExperimentalWarningApplyDuration,
		ExperimentalMemoryMlock:                  cfg.ExperimentalMemoryMlock,
		ExperimentalTxnModeWriteWithSharedBuffer: cfg.ExperimentalTxnModeWriteWithSharedBuffer,
		ExperimentalBootstrapDefragThresholdMegabytes: cfg.ExperimentalBootstrapDefragThresholdMegabytes,
		V2Deprecation: cfg.V2DeprecationEffective(),
	}

	var server *etcdserver.EtcdServer
	if server, err = etcdserver.NewServer(srvcfg); err != nil {
		return nil, err
	}
	if memberInitialized {
		if err = server.CheckInitialHashKV(); err != nil {
			log.Print("checkInitialHashKV failed", zap.Error(err))
			server.Cleanup()
			server = nil
			return nil, err
		}
	}
	server.Start()
	return server, nil
}

// NewEtcdServer 新建etcd服务，isJoin为true表示加入一个新的集群
func NewEtcdServer(name string, clients []string, peers []string, clusters map[string][]string, isJoin bool) (*etcdserver.EtcdServer, error) {
	var err error
	cfg := embed.NewConfig()
	cfg.Name = name
	if isJoin {
		cfg.ClusterState = embed.ClusterStateFlagExisting
	}
	cfg.InitialClusterToken = defaultClusterName
	cfg.Dir = defaultDir
	cfg.AutoCompactionMode = defaultAutoCompactionMode
	cfg.LCUrls, err = parseAndCheckURLs(clients)
	if err != nil {
		return nil, err
	}
	cfg.LPUrls, err = parseAndCheckURLs(peers)
	if err != nil {
		return nil, err
	}
	cfg.ACUrls, err = parseAndCheckURLs(clients)
	if err != nil {
		return nil, err
	}
	cfg.APUrls, err = parseAndCheckURLs(peers)
	if err != nil {
		return nil, err
	}
	cfg.InitialCluster = initialClusterString(clusters)
	return startEtcdServer(cfg)
}

// parseAndCheckURLs parses list of strings to url.URL objects.
func parseAndCheckURLs(urlStrings []string) ([]url.URL, error) {
	urls := make([]url.URL, len(urlStrings))
	for i, urlString := range urlStrings {
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			return nil, fmt.Errorf(" %s: %v", urlString, err)
		}
		urls[i] = *parsedURL
	}
	return urls, nil
}
func isMemberInitialized(cfg *embed.Config) bool {
	waldir := cfg.WalDir
	if waldir == "" {
		waldir = filepath.Join(cfg.Dir, "member", "wal")
	}
	return wal.Exist(waldir)
}
func initialClusterString(clusters map[string][]string) string {
	ss := make([]string, 0)
	for name, urls := range clusters {
		for _, u := range urls {
			ss = append(ss, fmt.Sprintf("%s=%s", name, u))
		}
	}
	return strings.Join(ss, ",")

}
