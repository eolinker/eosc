package raft_service

func (s *Service) AddEventHandler(h IRaftEventHandler) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.eventHandlers = append(s.eventHandlers, h)
}

func (s *Service) AddCommitEventHandler(h ICommitEventHandler) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.commitEventHandlers = append(s.commitEventHandlers, h)
}

func (s *Service) callbackEvent(event string) {
	hs := s.eventHandlers
	for _, h := range hs {
		h(event)
	}

}
func (s *Service) callCommitEvent(namespace, cmd string) {
	hs := s.commitEventHandlers
	for _, h := range hs {
		h(namespace, cmd)
	}

}
