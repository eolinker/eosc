package raft

import "go.etcd.io/etcd/raft/v3"

// 业务处理，根据实际需求更改service，service是外层的业务对象

type IRaftService interface {
	// Commit 节点commit信息前的处理
	Commit(command string, namespace string, key string, data []byte) (err error)

	// GetInit 集群初始化时的将service缓存中的信息进行打包处理,只会在切换集群模式的时候调用一次
	GetInit() (data []byte, err error)

	// ResetSnap 读取快照，用于恢复service数据
	ResetSnap(data []byte, isInit bool) (err error)

	// GetSnapshot 生成快照，用于快照文件的生成
	GetSnapshot() (data []byte, err error)
}
type IRaftStateHandler interface {
	SetState(stateType raft.StateType)
}

type IRaftStateHandlers []IRaftStateHandler

func (hs IRaftStateHandlers) SetState(stateType raft.StateType) {
	for _, h := range hs {
		h.SetState(stateType)
	}
}

type emptyIRaftStateHandlers struct {
}

func (e *emptyIRaftStateHandlers) SetState(stateType raft.StateType) {

}

func CreateRaftStateHandlers(handlers ...IRaftStateHandler) IRaftStateHandler {
	if len(handlers) > 0 {
		return IRaftStateHandlers(handlers)
	}
	return new(emptyIRaftStateHandlers)
}
