package raft_service

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var (
	ErrInvalidNamespace     = errors.New("invalid namespace")
	ErrInvalidCommitHandler = errors.New("invalid commit handler")
	ErrInvalidCommand       = errors.New("invalid command")
)

const (
	CommandInit     = "init"
	SystemNamespace = "__system"
)

type Service struct {
	handlers eosc.IUntyped
	raftNode raft.IRaftSender
	//processHandlers eosc.IUntyped
}

func (s *Service) Send(namespace, cmd string, body []byte) error {

	data, err := encodeCmd(namespace, cmd, body)
	if err != nil {
		return err
	}
	return s.raftNode.Send(data)
}

func (s *Service) CommitHandler(data []byte) error {
	cmd, err := unMarshalCmd(data)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return s.commitHandler(cmd.Namespace, cmd.Cmd, cmd.Body)

}

func (s *Service) ProcessDataHandler(body []byte) (data []byte, err error) {
	cmd, err := unMarshalCmd(data)
	if err != nil {
		return nil, err
	}
	return s.ProcessHandler(cmd.Namespace, cmd.Cmd, cmd.Body)
}

func (s *Service) SetRaft(raft raft.IRaftSender) {
	s.raftNode = raft
}

func (s *Service) ProcessHandler(namespace string, command string, processData []byte) ([]byte, error) {
	v, has := s.handlers.Get(namespace)
	if !has {
		return nil, ErrInvalidNamespace
	}
	f, ok := v.(IRaftServiceHandler)
	if !ok {
		return nil, ErrInvalidCommitHandler
	}
	body, err := f.ProcessHandler(command, processData)
	if err != nil {
		return nil, err
	}
	return encodeCmd(namespace, command, body)

}

//func RegisterHandlers(s *Service, handlers ...ICreateHandler) {
//	if handlers != nil {
//		for _, cf := range handlers {
//			h, ok := cf.(ICreateHandler)
//			if !ok {
//				continue
//			}
//			s.SetHandler(h.Namespace(), h.Handler())
//		}
//	}
//}

func (s *Service) SetHandlers(handlers ...ICreateHandler) {
	if handlers != nil {
		for _, cf := range handlers {
			h, ok := cf.(ICreateHandler)
			if !ok {
				continue
			}
			s.SetHandler(h.Namespace(), h.Handler())
		}
	}
}

func NewService() *Service {
	s := &Service{
		handlers: eosc.NewUntyped(),
	}
	return s
}

func (s *Service) SetHandler(namespace string, handler IRaftServiceHandler) {
	s.handlers.Set(namespace, handler)
}

func (s *Service) commitHandler(namespace string, cmd string, data []byte) error {
	if namespace == SystemNamespace {
		switch cmd {
		case CommandInit:
			return s.ResetSnap(data)
		}
		return ErrInvalidCommand
	}

	v, has := s.handlers.Get(namespace)
	if !has {
		return ErrInvalidNamespace
	}
	f, ok := v.(IRaftServiceHandler)
	if !ok {
		return ErrInvalidCommitHandler
	}

	return f.CommitHandler(cmd, data)

}

func (s *Service) GetInit() ([]byte, error) {

	data, err := s.GetSnapshot()
	if err != nil {
		return nil, err
	}
	cmd := &Commend{
		Namespace: SystemNamespace,
		Cmd:       CommandInit,
		Body:      data,
	}
	return proto.Marshal(cmd)

}

func (s *Service) ResetSnap(data []byte) error {
	snaps := make(map[string]string)
	err := json.Unmarshal(data, &snaps)
	if err != nil {
		return err
	}
	for namespace, value := range snaps {
		handler, has := s.handlers.Get(namespace)
		if !has {
			log.Warnf("reset snap %s:%w", namespace, ErrInvalidNamespace)
			continue
		}
		h, ok := handler.(IRaftServiceHandler)
		if !ok {
			log.Warnf("reset snap %s:%w", namespace, ErrInvalidCommitHandler)
			continue
		}
		d, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Errorf("reset snap %s:%w", namespace, err)
			continue
		}
		err = h.ResetHandler(d)
		if err != nil {
			log.Errorf("reset snap %s:%w", namespace, err)
			continue
		}
	}
	return nil
}

func (s *Service) GetSnapshot() ([]byte, error) {
	snapshots := make(map[string]string)
	for namespace, handler := range s.handlers.All() {
		h, ok := handler.(IRaftServiceHandler)
		if !ok {
			continue
		}
		snapshots[namespace] = base64.StdEncoding.EncodeToString(h.Snapshot())
	}
	return json.Marshal(snapshots)
}
