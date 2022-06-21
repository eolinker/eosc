package etcdRaft

import (
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"net/http"
	"os"
	"strings"
)

func (e *EtcdServer) cleanWalFile() error {
	if fileutil.Exist(defaultDir) {
		err := os.RemoveAll(defaultDir)
		if err != nil {
			return fmt.Errorf("eosc: cannot remove old dir for wal (%w)", err)
		}
	}
	return nil
}
func GenKey(prefix string, namespace string, key string) string {
	return fmt.Sprintf("%s/%s/%s", prefix, namespace, key)
}
func SpiltKey(prefix string, k string) (namespace string, key string, err error) {
	cleanKey := strings.TrimPrefix(k, prefix)
	res := strings.Split(cleanKey, "/")
	l := len(res)
	if l < 2 {
		return "", "", fmt.Errorf("invalid key : %s", cleanKey)
	}
	return res[l-2], res[l-1], nil
}

func GetRealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.RemoteAddr
	}
	return realIP
}
