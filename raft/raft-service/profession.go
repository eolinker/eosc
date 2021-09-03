package raft_service

import "github.com/eolinker/eosc"

type Profession struct {
}

func (p *Profession) ProcessHandler(propose []byte) (string, []byte, error) {
	return eosc.SpaceProfession, propose, nil
}

func (p *Profession) CommitHandler(data []byte) error {
	return nil
}
