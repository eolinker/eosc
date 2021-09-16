package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/store"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	store       eosc.IStore
	fileName    string
	professions []interface{}
}

func (p *Professions) ResetHandler(data []byte) error {
	return nil
}

func (p *Professions) CommitHandler(cmd string, data []byte) error {
	return nil
}

func (p *Professions) Snapshot() []byte {
	return nil
}

func (p *Professions) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	return nil, nil
}

func NewProfessions(fileName string) *Professions {
	return &Professions{
		store:    store.NewStore(),
		fileName: fileName,
	}
}
