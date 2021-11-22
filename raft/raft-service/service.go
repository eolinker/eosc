package raft_service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"

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
	appendDelay     = time.Second
	EventReset      = "reset"
	EventComplete   = "complete"
)

type Service struct {
	handlers eosc.IUntyped
	raftNode raft.IRaftSender

	commitHandler func(namespace, cmd string, body []byte) error

	timerAppend *time.Timer
	locker      sync.Mutex

	eventHandlers       []IRaftEventHandler
	commitEventHandlers []ICommitEventHandler
}

func NewService() *Service {
	s := &Service{
		handlers: eosc.NewUntyped(),
	}
	s.SetHandler(SystemNamespace, s)
	s.commitHandler = s.doCommit
	return s
}

func (s *Service) ResetHandler(data []byte) error {
	return errors.New("not support")
}

func (s *Service) Append(cmd string, data []byte) error {

	switch cmd {
	case CommandInit:
		snaps := make(map[string]string)
		err := json.Unmarshal(data, &snaps)
		if err != nil {
			return err
		}
		s.doResetSnap(snaps)
	}
	return ErrInvalidCommand
}

func (s *Service) Complete() error {

	return nil
}

func (s *Service) complete() error {

	for namespace, f := range s.handlers.All() {
		if namespace == SystemNamespace {
			continue
		}
		f.(IRaftServiceHandler).Complete()
	}

	return nil
}

func (s *Service) CommitHandler(cmd string, data []byte) error {
	switch cmd {
	case CommandInit:
		snaps := make(map[string]string)
		err := json.Unmarshal(data, &snaps)
		if err != nil {
			return err
		}
		s.doResetSnap(snaps)
		return nil
	}
	return ErrInvalidCommand
}

func (s *Service) Snapshot() []byte {
	return nil
}

func (s *Service) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {
	return nil, nil, errors.New("not support")
}

func (s *Service) Send(namespace, cmd string, body []byte) (interface{}, error) {

	data, err := encodeCmd(namespace, cmd, body)
	if err != nil {
		return nil, err
	}
	return s.raftNode.Send(data)
}

func (s *Service) Commit(data []byte) error {
	cmd, err := unMarshalCmd(data)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.commitHandler(cmd.Namespace, cmd.Cmd, cmd.Body)
}

func (s *Service) PreProcessData(inBody []byte) (object interface{}, data []byte, err error) {
	cmd, err := unMarshalCmd(inBody)
	if err != nil {
		return nil, nil, err
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.processHandler(cmd.Namespace, cmd.Cmd, cmd.Body)
}

func (s *Service) SetRaft(raft raft.IRaftSender) {
	s.raftNode = raft
}

func (s *Service) processHandler(namespace string, command string, processData []byte) (interface{}, []byte, error) {
	v, has := s.handlers.Get(namespace)
	if !has {
		return nil, nil, ErrInvalidNamespace
	}
	f, ok := v.(IRaftServiceHandler)
	if !ok {
		return nil, nil, ErrInvalidCommitHandler
	}
	body, obj, err := f.ProcessHandler(command, processData)
	if err != nil {
		return nil, nil, err
	}
	cmd, err := encodeCmd(namespace, command, body)
	if err != nil {
		return nil, nil, err
	}
	return obj, cmd, err

}

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

func (s *Service) SetHandler(namespace string, handler IRaftServiceHandler) {
	s.handlers.Set(namespace, handler)
}

func (s *Service) doCommit(namespace string, cmd string, data []byte) error {

	v, has := s.handlers.Get(namespace)
	if !has {
		return ErrInvalidNamespace
	}
	f, ok := v.(IRaftServiceHandler)
	if !ok {
		return ErrInvalidCommitHandler
	}
	err := f.CommitHandler(cmd, data)
	if err != nil {
		return err
	}
	s.callCommitEvent(namespace, cmd)
	return nil
}

func (s *Service) doAppend(namespace string, cmd string, data []byte) error {

	s.timerAppend.Reset(appendDelay)
	v, has := s.handlers.Get(namespace)
	if !has {
		return ErrInvalidNamespace
	}
	f, ok := v.(IRaftServiceHandler)
	if !ok {
		return ErrInvalidCommitHandler
	}

	err := f.Append(cmd, data)

	return err
}

func (s *Service) GetInit() ([]byte, error) {

	data, err := s.GetSnapshot()
	if err != nil {
		return nil, err
	}

	return encodeCmd(SystemNamespace, CommandInit, data)

}
func (s *Service) doResetSnap(snaps map[string]string) {
	for namespace, value := range snaps {
		if namespace == SystemNamespace {
			continue
		}
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
}
func (s *Service) ResetSnap(data []byte) error {
	snaps := make(map[string]string)
	err := json.Unmarshal(data, &snaps)
	if err != nil {
		return err
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	s.doResetSnap(snaps)
	s.callbackEvent(EventReset)
	if s.timerAppend == nil {
		s.timerAppend = time.NewTimer(appendDelay)
		s.commitHandler = s.doAppend
		go s.appendSwitch(s.timerAppend)
	}

	return nil
}
func (s *Service) appendSwitch(t *time.Timer) {

	<-t.C
	s.locker.Lock()
	defer s.locker.Unlock()
	s.timerAppend = nil
	s.commitHandler = s.doCommit
	s.complete()
	s.callbackEvent(EventComplete)
}
func (s *Service) GetSnapshot() ([]byte, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	snapshots := make(map[string]string)
	for namespace, handler := range s.handlers.All() {
		if namespace == SystemNamespace {
			continue
		}
		h, ok := handler.(IRaftServiceHandler)
		if !ok {
			continue
		}

		snapshots[namespace] = base64.StdEncoding.EncodeToString(h.Snapshot())
	}
	return json.Marshal(snapshots)
}
