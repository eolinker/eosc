package raft_service

import (
	"google.golang.org/protobuf/proto"
)

type IService interface {
	Send(namespace, cmd string, body []byte) (interface{}, error)
}

type IRaftServiceHandler interface {
	// ResetHandler  初始化snapshot
	ResetHandler(data []byte) error
	// Append 追加需要初始化的日志
	Append(cmd string, data []byte) error
	// Complete 初始化完成
	Complete() error

	// CommitHandler 节点commit信息前的处理
	CommitHandler(cmd string, data []byte) error
	// Snapshot 获取快照
	Snapshot() []byte
	// ProcessHandler leader 节点propose信息前的处理
	ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error)
}

type IRaftEventHandler func(event string)
type ICommitEventHandler func(namespace, cmd string)

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
