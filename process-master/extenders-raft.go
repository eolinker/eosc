package process_master

type ExtenderRaft struct {
}

func (e *ExtenderRaft) ResetHandler(data []byte) error {
	panic("implement me")
}

func (e *ExtenderRaft) CommitHandler(cmd string, data []byte) error {
	panic("implement me")
}

func (e *ExtenderRaft) Snapshot() []byte {
	panic("implement me")
}

func (e *ExtenderRaft) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {
	panic("implement me")
}
