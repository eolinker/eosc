package raft_service

type ICommitHandler interface {
	// InitHandler 初始化日志操作
	ResetHandler(data []byte) error
	// CommitHandler 节点commit信息前的处理
	CommitHandler(data []byte) error

	Snapshot() []byte
}

type IProcessHandler interface {
	// ProcessHandler 节点propose信息前的处理
	ProcessHandler(propose []byte) (string, []byte, error)
}
