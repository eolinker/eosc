package raft_service

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"sync"
	"time"

	dispatcher "github.com/eolinker/eosc/common/dispatcher"
)

type Service struct {
	dispatcher.IDispatchCenter
	closeOnes sync.Once
	data      *dispatcher.Data
	initOps   []InitOp

	eventChan chan dispatcher.IEvent
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
		eventChan:       make(chan dispatcher.IEvent, 1),
	}
	go s.myEventLoop()
	return s
}

func (s *Service) Commit(command string, namespace string, key string, data []byte) (err error) {

	return s.sendEvent(namespace, command, key, data)
}

func (s *Service) myEventLoop() {
	t := time.NewTimer(time.Millisecond)
	t.Stop()
	defer t.Stop()
	isNeedInit := false
	for {
		select {
		case e, ok := <-s.eventChan:
			if !ok {
				return
			}
			s.data.DoEvent(e)
			switch e.Event() {
			case eosc.EventReset, eosc.EventInit:
				isNeedInit = true
				t.Reset(time.Millisecond)
			default:
				if isNeedInit {
					t.Reset(time.Millisecond)
				} else {
					s.IDispatchCenter.Send(e)
				}
			}
		case <-t.C:
			if isNeedInit {
				isNeedInit = false
				d := s.data.GET()
				s.IDispatchCenter.Send(&Event{
					namespace: "",
					cmd:       eosc.EventReset,
					key:       "",
					data:      nil,
					all:       d,
				})
			}
		}
	}
}
func (s *Service) sendEvent(namespace, cmd, key string, data []byte) error {
	event := &Event{
		namespace: namespace,
		cmd:       cmd,
		key:       key,
		data:      data,
	}
	s.eventChan <- event

	return nil
}

func (s *Service) GetInit() ([]byte, error) {
	data, err := s.GetSnapshot()
	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *Service) ResetSnap(data []byte, isInit bool) error {
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
	ev := &Event{
		namespace: "",
		cmd:       eosc.EventReset,
		key:       "",
		data:      nil,
		all:       all,
	}
	s.eventChan <- ev
	return nil
}

func (s *Service) GetSnapshot() ([]byte, error) {
	all := s.data.GET()
	data, err := json.Marshal(all)
	return data, err
}
