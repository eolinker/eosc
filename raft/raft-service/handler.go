package raft_service

import "google.golang.org/protobuf/proto"

type IService interface {
	Send(namespace, cmd string, body []byte) (interface{}, error)
}

type IRaftServiceHandler interface {
	// ResetHandler  初始化日志操作
	ResetHandler(data []byte) error
	// CommitHandler 节点commit信息前的处理
	CommitHandler(cmd string, data []byte) error

	DelayDone()

	CommitDelay(cmd string, data []byte) error
	// Snapshot 获取快照
	Snapshot() []byte
	// ProcessHandler 节点propose信息前的处理
	ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error)
}

func unMarshalCmd(data []byte) (*Commend, error) {
	cmd := new(Commend)
	err := proto.Unmarshal(data, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, err
}
func encodeCmd(namespace, command string, body []byte) ([]byte, error) {
	cmd := &Commend{
		Namespace: namespace,
		Cmd:       command,
		Body:      body,
	}
	data, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	return data, err
}
