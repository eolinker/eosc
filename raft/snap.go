package raft

import (
	"encoding/json"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-master/workers"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
)

// 读取快照文件
func (rc *Node) loadSnapshot() *raftpb.Snapshot {
	if wal.Exist(rc.waldir) {
		walSnaps, err := wal.ValidSnapshotEntries(rc.logger, rc.waldir)
		if err != nil {
			log.Fatalf("eosc: error listing snapshots (%v)", err)
		}
		// 获取最新的快照
		snapshot, err := rc.snapshotter.LoadNewestAvailable(walSnaps)
		if err != nil && err != snap.ErrNoSnapshot {
			log.Fatalf("eosc: error loading snapshot (%v)", err)
		}
		return snapshot
	}
	return &raftpb.Snapshot{}
}

// ReadSnap 读取快照内容到service
func (rc *Node) ReadSnap(snapshotter *snap.Snapshotter, init bool) error {
	// 读取快照的所有内容
	snapshot, err := snapshotter.Load()
	if err != nil {
		// 快照不存在
		if err != snap.ErrNoSnapshot {
			return err
		}
		if init {
			log.Infof("reset snapshot")
			snaps := map[string]string{
				workers.SpaceWorker: "",
			}
			data, _ := json.Marshal(snaps)
			return rc.service.ResetSnap(data)
		}

	}

	// 快照不为空的话写进service
	if snapshot != nil {
		// 将快照内容缓存到service中
		log.Infof("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		err = rc.service.ResetSnap(snapshot.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

// 保存快照文件
func (rc *Node) saveSnap(snap raftpb.Snapshot) error {
	walSnap := walpb.Snapshot{
		Index:     snap.Metadata.Index,
		Term:      snap.Metadata.Term,
		ConfState: &snap.Metadata.ConfState,
	}
	// save the snapshot file before writing the snapshot to the wal.
	// This makes it possible for the snapshot file to become orphaned, but prevents
	// a WAL snapshot entry from having no corresponding snapshot file.
	if err := rc.snapshotter.SaveSnap(snap); err != nil {
		return err
	}
	if err := rc.wal.SaveSnapshot(walSnap); err != nil {
		return err
	}
	// 压缩日志
	return rc.wal.ReleaseLockTo(snap.Metadata.Index)
}

// service 从快照中重取数据
func (rc *Node) publishSnapshot(snapshotToSave raftpb.Snapshot) {
	if raft.IsEmptySnap(snapshotToSave) {
		return
	}
	log.Infof("publishing snapshot at index %d", rc.snapshotIndex)
	defer log.Infof("finished publishing snapshot at index %d", rc.snapshotIndex)

	if snapshotToSave.Metadata.Index <= rc.appliedIndex {
		log.Fatalf("snapshot index [%d] should > progress.appliedIndex [%d]", snapshotToSave.Metadata.Index, rc.appliedIndex)
	}
	err := rc.ReadSnap(rc.snapshotter, false)
	if err != nil {
		log.Info("read snap from snap shotter error:", err)
	}
	rc.confState = snapshotToSave.Metadata.ConfState
	rc.snapshotIndex = snapshotToSave.Metadata.Index
	rc.appliedIndex = snapshotToSave.Metadata.Index
}

// 保存现有快照
func (rc *Node) maybeTriggerSnapshot() {
	// 还不到保存快照的数里
	if rc.appliedIndex-rc.snapshotIndex <= rc.snapCount {
		return
	}
	log.Infof("start snapshot [applied index: %d | last snapshot index: %d]", rc.appliedIndex, rc.snapshotIndex)

	// 获取service中的信息
	data, err := rc.service.GetSnapshot()
	if err != nil {
		log.Panic(err)
	}
	// 利用现有信息生成要保存的快照信息
	snapContent, err := rc.raftStorage.CreateSnapshot(rc.appliedIndex, &rc.confState, data)
	if err != nil {
		log.Panic(err)
	}
	if err = rc.saveSnap(snapContent); err != nil {
		log.Panic(err)
	}
	// 暂时先不做日志的压缩处理，有bug如下：
	// 已有集群中的节点生成快照后，后续新节点加入集群时无法成功同步
	// 修复处理：集群节点配置更新的时候也更新raftStorage中已有快照的confState
	compactIndex := uint64(1)
	if rc.appliedIndex > snapshotCatchUpEntriesN {
		compactIndex = rc.appliedIndex - snapshotCatchUpEntriesN
	}
	if err = rc.raftStorage.Compact(compactIndex); err != nil {
		log.Panic(err)
	}
	log.Infof("compacted log at index %d", compactIndex)
	rc.snapshotIndex = rc.appliedIndex
}
