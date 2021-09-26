package raft

import (
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"os"

	"go.etcd.io/etcd/server/v3/wal"

	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
)

// 从现有文件中读取日志
func (rc *Node) replayWAL() *wal.WAL {
	// 先获取现有快照
	snapshot := rc.loadSnapshot()
	// 再获取现有日志
	w := rc.openWAL(snapshot)
	_, st, ents, err := w.ReadAll()
	if err != nil {
		log.Fatalf("eosc: failed to read WAL (%v)", err)
	}

	// 节点日志缓存初始化
	rc.raftStorage = NewMemoryStorage()
	if snapshot != nil {

		err = rc.raftStorage.ApplySnapshot(*snapshot)
		if err != nil {
			log.Infof("eosc: failed to apply snapshot for raftStorage (%v)", err)
		}
	}
	err = rc.raftStorage.SetHardState(st)
	if err != nil {
		log.Infof("eosc: failed to set hardState for raftStorage (%v)", err)
	}
	// append to storage so raft starts at the right place in log
	err = rc.raftStorage.Append(ents)
	if err != nil {
		log.Infof("eosc: failed to append ents for raftStorage (%v)", err)
	}
	return w
}

// 读取(创建)wal日志文件
func (rc *Node) openWAL(snapshot *raftpb.Snapshot) *wal.WAL {
	if !wal.Exist(rc.waldir) {
		// 创建本地文件
		if err := os.Mkdir(rc.waldir, 0750); err != nil {
			log.Fatalf("eosc: cannot create dir for wal (%v)", err)
		}
		// 创建wal日志对象
		w, err := wal.Create(zap.NewExample(), rc.waldir, nil)
		if err != nil {
			log.Fatalf("eosc: create wal error (%v)", err)
		}
		w.Close()
	}
	// 该结构用于日志对象定位快照索引，方便日志读取
	walsnap := walpb.Snapshot{}
	// 获取目前快照的最新记录
	if snapshot != nil {
		walsnap.Index, walsnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}
	log.Infof("loading WAL at term %d and index %d", walsnap.Term, walsnap.Index)
	// 开启日志
	w, err := wal.Open(zap.NewExample(), rc.waldir, walsnap)
	if err != nil {
		log.Fatalf("eosc: error loading wal (%v)", err)
	}
	return w
}

func (rc *Node) removeWalFile() error {
	if fileutil.Exist(rc.waldir) {
		err := os.RemoveAll(rc.waldir)
		if err != nil {
			return fmt.Errorf("eosc: cannot remove old dir for wal (%w)", err)
		}
	}
	if fileutil.Exist(rc.snapdir) {
		err := os.RemoveAll(rc.snapdir)
		if err != nil {
			return fmt.Errorf("eosc: cannot remove old dir for snap (%w)", err)
		}
	}
	return nil
}
