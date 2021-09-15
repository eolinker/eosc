package raft_service

type ICommitHandler interface {
	// ResetHandler  初始化日志操作
	ResetHandler(data []byte) error
	// CommitHandler 节点commit信息前的处理
	CommitHandler(data []byte) error

	Snapshot() []byte
}

type IProcessHandler interface {
	// ProcessHandler 节点propose信息前的处理
	ProcessHandler(propose interface{}) ([]byte, error)
}

type IRaftServiceHandler interface {
	ICommitHandler
	IProcessHandler
}

type RaftServiceHandler struct {
	ICommitHandler
	IProcessHandler
}

func NewRaftServiceHandler(ICommitHandler ICommitHandler, IProcessHandler IProcessHandler) *RaftServiceHandler {
	return &RaftServiceHandler{ICommitHandler: ICommitHandler, IProcessHandler: IProcessHandler}
}
