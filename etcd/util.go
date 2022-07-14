package etcd

import (
	"fmt"
	"github.com/eolinker/eosc/env"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"os"
	"path/filepath"
)

func (s *_Server) cleanWalFile() error {
	dir:=filepath.Join(env.DataDir(),"member")

	if fileutil.Exist(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			return fmt.Errorf("eosc: cannot remove old dir for wal (%w)", err)
		}
	}
	return nil
}
