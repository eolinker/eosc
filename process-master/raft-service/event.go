package raft_service

type Event struct {
	namespace string
	cmd       string
	key       string
	data      []byte
	all       map[string]map[string][]byte
}

func (e *Event) All() map[string]map[string][]byte {
	return e.all
}

func (e *Event) Namespace() string {
	return e.namespace
}

func (e *Event) Event() string {
	return e.cmd
}

func (e *Event) Key() string {
	//TODO implement me
	return e.key
}

func (e *Event) Data() []byte {
	return e.data
}
