package raft_service

import (
	"github.com/eolinker/eosc"
	dispatcher "github.com/eolinker/eosc/common/dispatcher"
	"github.com/eolinker/eosc/etcd"
	"strings"
	"sync"
)

type Service struct {
	dispatcher.IDispatchCenter
	closeOnes sync.Once

	initOps []InitOp
}

func (s *Service) Put(key string, value []byte) error {
	namespace, k := readKeys(key)
	return s.sendEvent(namespace, eosc.EventSet, k, value)
}

func (s *Service) Delete(key string) error {
	namespace, k := readKeys(key)
	return s.sendEvent(namespace, eosc.EventDel, k, nil)
}

func (s *Service) Reset(values []*etcd.KValue) {
	vs := make(map[string]map[string][]byte)

	for _, v := range values {
		namespace, k := readKeys(string(v.Key))

		if _, has := vs[namespace]; !has {
			vs[namespace] = make(map[string][]byte)
		}
		vs[namespace][k] = v.Value
	}
	for _, oh := range s.initOps {
		vs = oh(vs)
	}
	event := &Event{
		namespace: "",
		cmd:       eosc.EventReset,
		key:       "",
		data:      nil,
		all:       vs,
	}
	s.IDispatchCenter.Send(event)
}
func readKeys(keys string) (namespace, key string) {
	i := strings.LastIndex(keys, "/")
	if i > 0 {
		namespace = strings.TrimPrefix(keys[:i], "/")
		key = keys[i+1:]
	} else {
		namespace = ""
		key = strings.TrimPrefix(keys, "/")
	}
	return
}

type InitOp func(map[string]map[string][]byte) map[string]map[string][]byte

func NewService(ops ...InitOp) *Service {

	s := &Service{
		IDispatchCenter: dispatcher.NewDataDispatchCenter(),
		initOps:         ops,
	}

	return s
}

func (s *Service) sendEvent(namespace, cmd, key string, data []byte) error {
	event := &Event{
		namespace: namespace,
		cmd:       cmd,
		key:       key,
		data:      data,
	}
	s.IDispatchCenter.Send(event)
	return nil
}
