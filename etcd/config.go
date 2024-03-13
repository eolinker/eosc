// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/config"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3compactor"
	"go.etcd.io/etcd/server/v3/wal"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultMaxSnapshots = 5
	DefaultMaxWALs      = 5

	DefaultWarningApplyDuration = 100 * time.Millisecond

	DefaultDowngradeCheckTime = 5 * time.Second

	DefaultStrictReconfigCheck = true
	// DefaultEnableV2 is the default value for "--enable-v2" flag.
	// v2 API is disabled by default.
	DefaultEnableV2 = false
)

type Config struct {
	PeerAdvertiseUrls    []string
	ClientAdvertiseUrls  []string
	GatewayAdvertiseUrls []string
	DataDir              string
}

func (c *Config) CreatePeerUrl() (types.URLs, types.URLs, error) {

	peerUrl, err := parseAndCheckURLs(c.PeerAdvertiseUrls)
	if err != nil {
		return nil, nil, err
	}
	clientUrl, err := parseAndCheckURLs(c.ClientAdvertiseUrls)
	if err != nil {
		return nil, nil, err
	}

	return peerUrl, clientUrl, nil
}

func (s *_Server) etcdServerConfig() config.ServerConfig {

	dataDir := s.config.DataDir

	peerUrl, clientUrl, err := s.config.CreatePeerUrl()
	if err != nil {
		log.Fatal(err)
	}

	name, InitialCluster, isNew := s.readCluster(peerUrl)
	os.Setenv("node_id", name)
	var urlsmap types.URLsMap
	var token string
	if !wal.Exist(filepath.Join(dataDir, "member", "wal")) {
		urlsmap, err = types.NewURLsMap(InitialCluster)
		token = "APINTO_CLUSTER"
	}

	srvCfg := config.ServerConfig{

		ClientURLs:                               clientUrl,
		PeerURLs:                                 peerUrl,
		DataDir:                                  dataDir,
		DedicatedWALDir:                          "",
		MaxWALFiles:                              DefaultMaxWALs,
		InitialPeerURLsMap:                       urlsmap,
		InitialClusterToken:                      token,
		NewCluster:                               isNew,
		AutoCompactionRetention:                  time.Duration(10),
		AutoCompactionMode:                       v3compactor.ModeRevision,
		QuotaBackendBytes:                        quotaBackendBytes,
		BackendFreelistType:                      bolt.FreelistMapType,
		MaxTxnOps:                                maxTxnOps,
		MaxRequestBytes:                          maxRequestBytes,
		SocketOpts:                               transport.SocketOpts{},
		Logger:                                   zap.New(NewLogger()),
		DowngradeCheckTime:                       DefaultDowngradeCheckTime,
		WarningApplyDuration:                     DefaultWarningApplyDuration,
		MaxSnapFiles:                             DefaultMaxSnapshots,
		Name:                                     name,
		SnapshotCount:                            etcdserver.DefaultSnapshotCount,
		SnapshotCatchUpEntries:                   etcdserver.DefaultSnapshotCatchUpEntries,
		TickMs:                                   100,
		ElectionTicks:                            10,
		InitialElectionTickAdvance:               true,
		StrictReconfigCheck:                      DefaultStrictReconfigCheck,
		EnableGRPCGateway:                        true,
		CORS:                                     map[string]struct{}{"*": {}},
		HostWhitelist:                            map[string]struct{}{"*": {}},
		AuthToken:                                "simple",
		BcryptCost:                               uint(bcrypt.DefaultCost),
		TokenTTL:                                 300,
		PreVote:                                  true,
		ExperimentalMemoryMlock:                  false,
		ExperimentalTxnModeWriteWithSharedBuffer: true,
		V2Deprecation:                            config.V2_DEPR_DEFAULT,
	}
	return srvCfg
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

func initialClusterString(clusters map[string][]string) string {
	ss := make([]string, 0)
	for name, urls := range clusters {
		for _, u := range urls {
			ss = append(ss, fmt.Sprintf("%s=%s", name, u))
		}
	}
	return strings.Join(ss, ",")
}

func readClusterString(clusters string) map[string][]string {
	its := strings.Split(clusters, ",")
	vs := make(map[string][]string)
	for _, it := range its {
		f := strings.Split(it, "=")
		if len(f) != 2 {
			continue
		}
		vs[f[0]] = append(vs[f[0]], f[1])
	}
	return vs
}
func (s *_Server) resetCluster(InitialCluster string) {
	etcdInitPath := filepath.Join(s.config.DataDir, "cluster", "etcd.init")
	etcdConfig := env.NewConfig(etcdInitPath)
	etcdConfig.ReadFile(etcdInitPath)
	defer func() {
		err := etcdConfig.Save()
		if err != nil {
			log.Warn("write args file fail:", err)
			return
		}
		log.Info("write args file succeed!")
	}()
	etcdConfig.Set("cluster", InitialCluster)
}
func (s *_Server) updateCluster() {
	ctx, _ := s.requestContext()
	s.client.MemberList(ctx)
}
func (s *_Server) clearCluster() {
	s.resetCluster("")
}
func (s *_Server) readCluster(peerUrl types.URLs) (name, InitialCluster string, isNew bool) {
	etcdInitPath := filepath.Join(s.config.DataDir, "cluster", "etcd.init")
	etcdConfig := env.NewConfig(etcdInitPath)
	etcdConfig.ReadFile(etcdInitPath)
	defer func() {
		err := etcdConfig.Save()
		if err != nil {
			log.Warn("write args file fail:", err)
			return
		}
		log.Info("write args file succeed!")
	}()

	var has bool
	name, has = etcdConfig.Get("name")
	if !has {
		name = createUUID()
		etcdConfig.Set("name", name)
	}
	InitialCluster, has = etcdConfig.Get("cluster")

	if !has || InitialCluster == "" {
		isNew = true
		members := map[string][]string{
			name: peerUrl.StringSlice(),
		}
		InitialCluster = initialClusterString(members)
	} else {
		members := readClusterString(InitialCluster)
		members[name] = peerUrl.StringSlice()
		InitialCluster = initialClusterString(members)
	}
	etcdConfig.Set("cluster", InitialCluster)
	return name, InitialCluster, isNew
}

func createUUID() string {
	u := uuid.New()
	bs := make([]byte, 32)
	hex.Encode(bs, u[:])
	return string(bs)
}
