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
	eoscConfig "github.com/eolinker/eosc/config"
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
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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

var etcdInitPath = filepath.Join(env.DataDir(), "cluster", "etcd.init")

func CreatePeerUrl() (types.URLs, types.URLs, error) {
	c, err := eoscConfig.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	admin := c.Admin
	peerUrl, err := createPeerUrl(admin.Scheme, []int{admin.Listen}, []string{admin.IP})
	if err != nil {
		return nil, nil, err
	}
	clientUrl, err := createPeerUrl("eosc", c.Listen, nil)
	if err != nil {
		return nil, nil, err
	}
	return peerUrl, clientUrl, nil
}
func createPeerUrl(schema string, ports []int, ips []string) (types.URLs, error) {

	urls := make([]string, 0)
	for _, ip := range ips {
		if ip == "" || ip == "0.0.0.0" {
			return createPeerUrl(schema, ports, readAllIp())
		}
		for _, port := range ports {
			if schema != "" {
				urls = append(urls, fmt.Sprintf("%s://%s:%d", schema, ip, port))
			} else {
				urls = append(urls, fmt.Sprintf("%s:%d", ip, port))

			}
		}
	}
	if len(urls) == 0 {
		return createPeerUrl(schema, ports, readAllIp())
	}

	return parseAndCheckURLs(urls)
}
func readAllIp() []string {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Debug(err)
		os.Exit(1)
	}

	ips := make([]string, 0, len(addrs))
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	if len(ips) == 0 {

		log.Fatal("not find valid ip")
	}
	return ips

}
func etcdServerConfig() config.ServerConfig {

	dataDir := env.DataDir()

	peerUrl, clientUrl, err := CreatePeerUrl()
	if err != nil {
		log.Fatal(err)
	}

	name, InitialCluster, isNew := readCluster(peerUrl)
	var urlsmap types.URLsMap
	var token string
	if !wal.Exist(filepath.Join(dataDir, "member", "wal")) {
		urlsmap, err = types.NewURLsMap(InitialCluster)
		token = "APINTO_CLUSTER"
	}

	srvcfg := config.ServerConfig{

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
	return srvcfg
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
func resetCluster(InitialCluster string) {
	etcdConfig := env.NewConfig(etcdInitPath)
	etcdConfig.ReadFile(etcdInitPath)
	defer etcdConfig.Save()
	etcdConfig.Set("cluster", InitialCluster)
}

func clearCluster() {
	resetCluster("")
}
func readCluster(peerUrl types.URLs) (name, InitialCluster string, isNew bool) {

	etcdConfig := env.NewConfig(etcdInitPath)
	etcdConfig.ReadFile(etcdInitPath)
	defer etcdConfig.Save()

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
		etcdConfig.Set("cluster", InitialCluster)
	}
	return name, InitialCluster, isNew
}

func createUUID() string {
	u := uuid.New()
	bs := make([]byte, 32)
	hex.Encode(bs, u[:])
	return string(bs)
}
