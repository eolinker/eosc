package raft_service

import (
	"encoding/json"
	"sync"

	"github.com/eolinker/eosc"

	dispatcher "github.com/eolinker/eosc/common/dispatcher"
)

type Service struct {
	dispatcher.IDispatchCenter
	closeOnes sync.Once
	data      *dispatcher.Data
	initOps   []InitOp
}
type InitOp func(map[string]map[string][]byte) map[string]map[string][]byte

func NewService(ops ...InitOp) *Service {
	initData := make(map[string]map[string][]byte)
	for _, namespace := range eosc.Namespaces {
		initData[namespace] = make(map[string][]byte)
	}
	s := &Service{
		IDispatchCenter: dispatcher.NewDataDispatchCenter(),
		data:            dispatcher.NewMyData(initData),
		initOps:         ops,
	}

	return s
}

func (s *Service) Commit(command string, namespace string, key string, data []byte) (err error) {

	return s.sendEvent(namespace, command, key, data)
}

func (s *Service) sendEvent(namespace, cmd, key string, data []byte) error {
	event := &Event{
		namespace: namespace,
		cmd:       cmd,
		key:       key,
		data:      data,
	}

	s.data.DoEvent(event)
	//event.all = s.data.GET()
	s.IDispatchCenter.Send(event)

	return nil
}

func (s *Service) GetInit() ([]byte, error) {
	data, err := s.GetSnapshot()
	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *Service) ResetSnap(data []byte) error {
	all := make(map[string]map[string][]byte)
	if len(data) != 0 {
		err := json.Unmarshal(data, &all)
		if err != nil {
			return err
		}
	}
	for _, op := range s.initOps {
		all = op(all)
	}
	s.data.DoEvent(&Event{
		namespace: "",
		cmd:       eosc.EventReset,
		key:       "",
		data:      nil,
		all:       all,
	})
	return nil
}

func (s *Service) GetSnapshot() ([]byte, error) {
	all := s.data.GET()
	data, err := json.Marshal(all)
	return data, err
}
