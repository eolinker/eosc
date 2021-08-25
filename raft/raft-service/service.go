package raft_service

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var (
	ErrInvalidNamespace     = errors.New("invalid namespace")
	ErrInvalidCommitHandler = errors.New("invalid commit handler")
	ErrInvalidKey           = errors.New("invalid key")
	commandSet              = "set"
	commandDel              = "delete"
)

type Service struct {
	store           eosc.IStore
	commitHandlers  eosc.IUntyped
	processHandlers eosc.IUntyped
}

func NewService(store eosc.IStore) *Service {
	return &Service{
		store:           store,
		commitHandlers:  eosc.NewUntyped(),
		processHandlers: eosc.NewUntyped(),
	}
}

func (s *Service) CommitHandlerSet(namespace string, handler ICommitHandler) {
	s.commitHandlers.Set(namespace, handler)
}

func (s *Service) ProcessHandlerSet(namespace string, handler IProcessHandler) {
	s.processHandlers.Set(namespace, handler)
}

func (s *Service) CommitHandler(namespace string, data []byte) error {
	if namespace == "init" {
		return s.ResetSnap(data)
	}
	v, has := s.commitHandlers.Get(namespace)
	if !has {
		return ErrInvalidNamespace
	}
	f, ok := v.(ICommitHandler)
	if !ok {
		return ErrInvalidCommitHandler
	}
	return f.CommitHandler(data)
}

func (s *Service) ProcessHandler(namespace string, propose []byte) (string, []byte, error) {
	v, has := s.processHandlers.Get(namespace)
	if !has {
		return "", nil, ErrInvalidNamespace
	}
	f, ok := v.(IProcessHandler)
	if !ok {
		return "", nil, ErrInvalidCommitHandler
	}
	return f.ProcessHandler(propose)
}

func (s *Service) GetInit() (string, []byte, error) {
	data, err := s.GetSnapshot()
	return "init", data, err
}

func (s *Service) ResetSnap(data []byte) error {
	snaps := make(map[string]string)
	err := json.Unmarshal(data, &snaps)
	if err != nil {
		return err
	}
	for namespace, value := range snaps {
		handler, has := s.commitHandlers.Get(namespace)
		if !has {
			log.Warnf("reset snap %s:%w", namespace, ErrInvalidNamespace)
			continue
		}
		h, ok := handler.(ICommitHandler)
		if !ok {
			log.Warnf("reset snap %s:%w", namespace, ErrInvalidCommitHandler)
			continue
		}
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Errorf("reset snap %s:%w", namespace, err)
			continue
		}
		err = h.ResetHandler(data)
		if err != nil {
			log.Errorf("reset snap %s:%w", namespace, err)
			continue
		}
	}
	return nil
}

func (s *Service) GetSnapshot() ([]byte, error) {
	snapshots := make(map[string]string)
	for namespace, handler := range s.commitHandlers.All() {
		h, ok := handler.(ICommitHandler)
		if !ok {
			continue
		}
		snapshots[namespace] = base64.StdEncoding.EncodeToString(h.Snapshot())
	}
	return json.Marshal(snapshots)
}
