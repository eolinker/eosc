package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/store"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	store eosc.IStore
}

func NewProfessions() *Professions {
	return &Professions{
		store: store.NewStore(),
	}
}

func (p *Professions) ResetHandler(data []byte) error {
	return nil
}

func (p *Professions) CommitHandler(data []byte) error {
	panic("implement me")
}

func (p *Professions) Snapshot() []byte {
	panic("implement me")
}

func (p *Professions) ProcessHandler(propose []byte) (string, []byte, error) {
	panic("implement me")
}
