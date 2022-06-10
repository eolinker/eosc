package raft

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"google.golang.org/protobuf/proto"

	"github.com/eolinker/eosc/log"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

// publishEntries writes committed log entries to commit channel and returns
// whether all entries could be published.
// 日志提交处理
func (rc *Node) publishEntries(ents []raftpb.Entry) bool {
	if len(ents) == 0 {
		return true
	}
	for i := range ents {
		switch ents[i].Type {
		case raftpb.EntryNormal:
			if len(ents[i].Data) == 0 {
				// ignore empty messages
				continue
			}
			cmd := new(Command)
			err := proto.Unmarshal(ents[i].Data, cmd)
			if err != nil {
				log.Errorf("unmarshal raft message :%w", err)
				continue
			}
			type basic struct {
				ID string `json:"id"`
			}

			if cmd.Version == "" { // 兼容旧版本数据
				switch cmd.Cmd {
				case eosc.EventSet:
					b := new(basic)
					err = json.Unmarshal(cmd.Body, b)
					if err != nil {
						continue
					}
					cmd.Key = b.ID
				case eosc.EventDel:
					cmd.Key = string(cmd.Body)
				}

			}
			if cmd.Cmd == eosc.EventReset {
				err = rc.service.ResetSnap(cmd.Body, false)
			} else {
				err = rc.service.Commit(cmd.Cmd, cmd.Namespace, cmd.Key, cmd.Body)
			}
			if err != nil {
				log.Error(err)
			}
		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			cc.Unmarshal(ents[i].Data)

			rc.confState = *rc.node.ApplyConfChange(cc)
			// eosc: 修改快照里的confState
			rc.raftStorage.UpdateConState(&rc.confState)
			switch cc.Type {
			// 新增节点
			case raftpb.ConfChangeAddNode:
				if len(cc.Context) > 0 {
					var info NodeInfo
					err := json.Unmarshal(cc.Context, &info)
					if err != nil {
						log.Errorf("fail to publish Entries,error:%s", err.Error())
						continue
					}
					// transport不需要加自己
					if cc.NodeID != rc.nodeID {
						p := rc.transport.Get(types.ID(cc.NodeID))
						// 不存在才加进去
						if p == nil {
							rc.transport.AddPeer(types.ID(cc.NodeID), []string{info.Addr})
						}
					}
					_, ok := rc.peers.GetPeerByID(cc.NodeID)
					if !ok {
						// 不存在，新增
						rc.peers.SetPeer(cc.NodeID, &info)
					}
					if rc.peers.GetPeerNum() > 1 && !rc.join {
						rc.join = true
						rc.writeConfig()
					}
				}
			case raftpb.ConfChangeRemoveNode:
				if cc.NodeID == rc.nodeID {
					log.Info("current node has been removed from the cluster!")
					return false
				}
				p := rc.transport.Get(types.ID(cc.NodeID))
				if p != nil {
					// 存在才移除
					rc.transport.RemovePeer(types.ID(cc.NodeID))
				}
				_, ok := rc.peers.GetPeerByID(cc.NodeID)
				if ok {
					// 存在，减去
					rc.peers.DeletePeerByID(cc.NodeID)
				}
				if rc.peers.GetPeerNum() < 2 && rc.join {
					rc.join = false
					rc.writeConfig()
				}
			}
		}
	}
	// after commit, update appliedIndex
	rc.appliedIndex = ents[len(ents)-1].Index

	return true
}

// 处理本节点需要committed的日志
func (rc *Node) entriesToApply(ents []raftpb.Entry) (nents []raftpb.Entry) {
	if len(ents) == 0 {
		return ents
	}
	firstIdx := ents[0].Index
	if firstIdx > rc.appliedIndex+1 {
		log.Fatalf("first index of c.ommitted entry[%d] should <= progress.appliedIndex[%d]+1", firstIdx, rc.appliedIndex)
	}
	if rc.appliedIndex-firstIdx+1 < uint64(len(ents)) {
		nents = ents[rc.appliedIndex-firstIdx+1:]
	}
	return nents
}

// saveToStorage 内存信息存储
func (rc *Node) saveToStorage(rd raft.Ready) {
	var err error
	err = rc.wal.Save(rd.HardState, rd.Entries)
	if err != nil {
		log.Fatal("save wal error: ", err)
	}
	// 安装快照信息
	if !raft.IsEmptySnap(rd.Snapshot) {
		log.Info("begin save snap")
		err = rc.saveSnap(rd.Snapshot)
		if err != nil {
			log.Fatal("save snap error: ", err)
		}
		// 有快照信息则根据快照修改内存
		err = rc.raftStorage.ApplySnapshot(rd.Snapshot)
		if err != nil {
			log.Fatal("raftStorage apply snapshot error: ", err)
		}
		rc.publishSnapshot(rd.Snapshot)
	}
	// 日志信息写入节点缓存
	err = rc.raftStorage.Append(rd.Entries)
	if err != nil {
		log.Fatal("raftStorage append entries error: ", err)
	}
}
