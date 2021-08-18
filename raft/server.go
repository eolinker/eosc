package raft

// 业务处理，根据实际需求更改service，service是外层的业务对象

type IService interface {
	// CommitHandler 节点commit信息前的处理
	CommitHandler(cmd string, data []byte) (err error)

	// ProcessHandler 节点propose信息前的处理
	ProcessHandler(command string, propose []byte) (cmd string, data []byte, err error)

	// GetInit 集群初始化时的将service缓存中的信息进行打包处理,只会在切换集群模式的时候调用一次
	GetInit() (cmd string, data []byte, err error)

	// ResetSnap 读取快照，用于恢复service数据
	ResetSnap(data []byte) (err error)

	// GetSnapshot 生成快照，用于快照文件的生成
	GetSnapshot() (data []byte, err error)
}
