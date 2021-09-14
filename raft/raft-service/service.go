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
)

type Service struct {
	handlers  eosc.IUntyped
	//processHandlers eosc.IUntyped
}

func NewService(handlers ...ICreateHandler) *Service {
	s := &Service{
		handlers:  eosc.NewUntyped(),
	}

	if handlers != nil{
		for _, cf := range handlers {
			h, ok := cf.(ICreateHandler)
			if !ok {
				continue
			}
			s.SetHandler(h.Namespace(),h.Handler())
			//s.commitHandlerSet(namespace,h.CommitHandler() )
			//s.processHandlerSet(namespace, h.ProcessHandler())
		}
	}
	return s
}
//
//func (s *service) commitHandlerSet(namespace string, handler ICommitHandler) {
//	s.commitHandlers.Set(namespace, handler)
//}
//
//func (s *service) processHandlerSet(namespace string, handler IProcessHandler) {
//	s.processHandlers.Set(namespace, handler)
//}

func (s *Service)SetHandler(namespace string,handler IRaftServiceHandler)  {
	s.handlers.Set(namespace,handler)
}

func (s *Service) CommitHandler(namespace string, data []byte) error {
	if namespace == "init" {
		return s.ResetSnap(data)
	}
	v, has := s.handlers.Get(namespace)
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
	v, has := s.handlers.Get(namespace)
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
		handler, has := s.handlers.Get(namespace)
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
	for namespace, handler := range s.handlers.All() {
		h, ok := handler.(ICommitHandler)
		if !ok {
			continue
		}
		snapshots[namespace] = base64.StdEncoding.EncodeToString(h.Snapshot())
	}
	return json.Marshal(snapshots)
}
