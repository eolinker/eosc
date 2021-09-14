package raft_service

import "github.com/eolinker/eosc"

var ff = eosc.NewUntyped()

type CreateHandler func(s eosc.IStore) (ICommitHandler, IProcessHandler)

func Register(namespace string, createFunc CreateHandler) {
	ff.Set(namespace, createFunc)
}

func initHandler(s *Service) {
	for namespace, cf := range ff.All() {
		f, ok := cf.(CreateHandler)
		if !ok {
			continue
		}
		ch, ph := f(s.store)
		s.commitHandlerSet(namespace, ch)
		s.processHandlerSet(namespace, ph)
	}
}
