package raft_service

type IService interface {
	Send(namespace, cmd string, body interface{})
}

type IRaftServiceHandler interface {
	// ResetHandler  初始化日志操作
	ResetHandler(data []byte) error
	// CommitHandler 节点commit信息前的处理
	CommitHandler(cmd string, data []byte) error
	// Snapshot 获取快照
	Snapshot() []byte
	// ProcessHandler 节点propose信息前的处理
	ProcessHandler(cmd string, body interface{}) ([]byte, error)
}
